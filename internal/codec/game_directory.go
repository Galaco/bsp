// Package codec implements bounded, stateless container codecs.
package codec

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/galaco/bsp/lumps"
)

// The PC wire layout and compressed boundary rule are defined by dgamelump_t
// and CompressGameLump in ValveSoftware/source-sdk-2013 at
// b2705ba55b3b802b86ef2b2dbf97939c9d4fb685. The zero-offset terminal convention
// is also handled by TeamSpen210/srctools at
// cd7e84107815e24b424c180f46da538eb4986a31 (src/srctools/bsp.py).
const (
	gameEntrySize = 16
	gameMaxInt32  = int64(1<<31 - 1)
)

// GameDirectory borrows Raw and all entry payloads. Base is the absolute file
// offset of Raw. Keep Raw and Entries unchanged while using or relocating it.
type GameDirectory struct {
	Raw     []byte
	Base    int32
	Entries []lumps.GameLump
}

// ParseGameDirectory validates a bounded little-endian PC container at base.
// Absent bytes and a present zero-count directory remain distinct in Raw.
func ParseGameDirectory(raw []byte, base int32) (GameDirectory, error) {
	count, directoryEnd, err := gameDirectoryHeader(raw, base)
	if err != nil {
		return GameDirectory{}, err
	}
	directory := GameDirectory{Raw: raw, Base: base}
	if count == 0 {
		return directory, nil
	}
	directory.Entries = make([]lumps.GameLump, count)
	for index := range directory.Entries {
		entry := readGameEntry(raw, index)
		start, end, err := gameEntryRange(raw, base, count, directoryEnd, index, entry)
		if err != nil {
			return GameDirectory{}, fmt.Errorf("game entry %d: %w", index, err)
		}
		if entry.Offset != 0 {
			entry.StoredData = raw[int(start):int(end):int(end)]
		}
		directory.Entries[index] = entry
	}
	return directory, nil
}

// GameRelocation borrows an immutable source and records validated offset changes.
// Prepare it before allocating a complete output file.
type GameRelocation struct {
	raw   []byte
	count int
	delta int64
}

func PlanGameRelocation(directory GameDirectory, newBase int32) (GameRelocation, error) {
	count, _, err := validateGameDirectory(directory.Raw, directory.Base)
	if err != nil {
		return GameRelocation{}, err
	}
	if err := validateGameExtent(len(directory.Raw), newBase); err != nil {
		return GameRelocation{}, err
	}
	delta := int64(newBase) - int64(directory.Base)
	for index := 0; index < count; index++ {
		offset := readGameEntry(directory.Raw, index).Offset
		if offset != 0 && (int64(offset)+delta < 0 || int64(offset)+delta > gameMaxInt32) {
			return GameRelocation{}, fmt.Errorf("game entry %d: relocated offset is not representable", index)
		}
	}
	return GameRelocation{raw: directory.Raw, count: count, delta: delta}, nil
}

// CopyTo copies directly into its final destination, avoiding a staging payload.
func (plan GameRelocation) CopyTo(raw []byte) error {
	if len(raw) != len(plan.raw) {
		return io.ErrShortBuffer
	}
	copy(raw, plan.raw)
	for index := 0; index < plan.count; index++ {
		position := 4 + index*gameEntrySize + 8
		offset := int32(binary.LittleEndian.Uint32(raw[position : position+4]))
		if offset != 0 {
			binary.LittleEndian.PutUint32(raw[position:position+4], uint32(int64(offset)+plan.delta))
		}
	}
	return nil
}

// RelocateGameDirectory preserves all bytes except nonzero absolute offsets.
// Raw and Base are authoritative; validation precedes one output allocation.
func RelocateGameDirectory(directory GameDirectory, newBase int32) ([]byte, error) {
	plan, err := PlanGameRelocation(directory, newBase)
	if err != nil {
		return nil, err
	}
	if len(directory.Raw) == 0 {
		return nil, nil
	}
	raw := make([]byte, len(directory.Raw))
	if err := plan.CopyTo(raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// EncodeGameDirectory emits a canonical directory followed by stored payloads,
// without per-payload length words. Uncompressed lengths are derived from data;
// compressed logical lengths are retained. Input offsets are replaced, except
// for zero-offset empty entries and the terminal sentinel's zero convention.
// No caller data is modified, including on validation failure.
func EncodeGameDirectory(entries []lumps.GameLump, base int32) ([]byte, error) {
	return encodeGameDirectory(entries, base, nil, nil)
}

func encodeGameDirectory(entries []lumps.GameLump, base int32, leading []byte, gaps []gameGap) ([]byte, error) {
	if base < 0 {
		return nil, fmt.Errorf("game directory: negative base %d", base)
	}
	count := int64(len(entries))
	if count > (gameMaxInt32-4)/gameEntrySize {
		return nil, fmt.Errorf("game directory: count %d exceeds signed 32-bit directory size", count)
	}
	total := int64(4) + count*gameEntrySize
	total += int64(len(leading))
	for index, entry := range entries {
		if entry.Offset < 0 {
			return nil, fmt.Errorf("game entry %d: negative offset", index)
		}
		if isGameSentinel(entry) && index == len(entries)-1 {
			if len(entry.StoredData) != 0 {
				return nil, fmt.Errorf("game entry %d: terminal boundary has payload", index)
			}
		} else if entry.Flags&1 != 0 {
			if entry.LogicalLength < 0 {
				return nil, fmt.Errorf("game entry %d: negative logical length", index)
			}
			if index == len(entries)-1 {
				return nil, fmt.Errorf("game entry %d: compressed payload lacks following boundary", index)
			}
		}
		storedLength := int64(len(entry.StoredData))
		total += (4 - (int64(base)+total)%4) % 4
		if storedLength > gameMaxInt32 || total+storedLength > gameMaxInt32 {
			return nil, fmt.Errorf("game entry %d: stored length or container size exceeds signed 32-bit range", index)
		}
		total += storedLength
		if gaps != nil {
			total += int64(len(gaps[index].tail)) + int64(len(gaps[index].after))
			if total > gameMaxInt32 {
				return nil, fmt.Errorf("game entry %d: retained gaps exceed signed 32-bit range", index)
			}
		}
	}
	if int64(base)+total > gameMaxInt32 {
		return nil, fmt.Errorf("game directory: absolute end exceeds signed 32-bit range")
	}

	raw := make([]byte, int(total))
	binary.LittleEndian.PutUint32(raw[:4], uint32(count))
	cursor := 4 + len(entries)*gameEntrySize
	copy(raw[cursor:], leading)
	cursor += len(leading)
	for index, entry := range entries {
		position := 4 + index*gameEntrySize
		cursor += int((4 - (int64(base)+int64(cursor))%4) % 4)
		offset := int32(int64(base) + int64(cursor))
		logicalLength := entry.LogicalLength
		sentinel := index == len(entries)-1 && isGameSentinel(entry)
		if sentinel {
			if entry.Offset == 0 {
				offset = 0
			}
		} else if entry.Flags&1 == 0 {
			logicalLength = int32(len(entry.StoredData))
			// A real entry following compressed data must provide its boundary,
			// even when that entry's own payload is empty.
			if len(entry.StoredData) == 0 && entry.Offset == 0 && (index == 0 || entries[index-1].Flags&1 == 0) {
				offset = 0
			}
		}
		binary.LittleEndian.PutUint32(raw[position:position+4], entry.ID)
		binary.LittleEndian.PutUint16(raw[position+4:position+6], entry.Flags)
		binary.LittleEndian.PutUint16(raw[position+6:position+8], entry.Version)
		binary.LittleEndian.PutUint32(raw[position+8:position+12], uint32(offset))
		binary.LittleEndian.PutUint32(raw[position+12:position+16], uint32(logicalLength))
		copy(raw[cursor:], entry.StoredData)
		cursor += len(entry.StoredData)
		if gaps != nil {
			copy(raw[cursor:], gaps[index].tail)
			cursor += len(gaps[index].tail)
			copy(raw[cursor:], gaps[index].after)
			cursor += len(gaps[index].after)
		}
	}
	return raw, nil
}

func validateGameExtent(length int, base int32) error {
	if base < 0 {
		return fmt.Errorf("game directory: negative base %d", base)
	}
	if int64(length) > gameMaxInt32 || int64(base)+int64(length) > gameMaxInt32 {
		return fmt.Errorf("game directory: container extent exceeds signed 32-bit range")
	}
	return nil
}

func gameDirectoryHeader(raw []byte, base int32) (int, int64, error) {
	if err := validateGameExtent(len(raw), base); err != nil {
		return 0, 0, err
	}
	if len(raw) == 0 {
		return 0, 0, nil
	}
	if len(raw) < 4 {
		return 0, 0, fmt.Errorf("game directory count: %w", io.ErrUnexpectedEOF)
	}
	count := int32(binary.LittleEndian.Uint32(raw[:4]))
	if count < 0 {
		return 0, 0, fmt.Errorf("game directory: negative count %d", count)
	}
	directoryEnd := int64(4) + int64(count)*gameEntrySize
	if directoryEnd > gameMaxInt32 {
		return 0, 0, fmt.Errorf("game directory: count %d exceeds signed 32-bit directory size", count)
	}
	if directoryEnd > int64(len(raw)) {
		return 0, 0, fmt.Errorf("game directory with %d entries: %w", count, io.ErrUnexpectedEOF)
	}
	return int(count), directoryEnd, nil
}

func validateGameDirectory(raw []byte, base int32) (int, int64, error) {
	count, directoryEnd, err := gameDirectoryHeader(raw, base)
	if err != nil {
		return 0, 0, err
	}
	for index := 0; index < count; index++ {
		entry := readGameEntry(raw, index)
		if _, _, err := gameEntryRange(raw, base, count, directoryEnd, index, entry); err != nil {
			return 0, 0, fmt.Errorf("game entry %d: %w", index, err)
		}
	}
	return count, directoryEnd, nil
}

func readGameEntry(raw []byte, index int) lumps.GameLump {
	position := 4 + index*gameEntrySize
	return lumps.GameLump{
		ID:            binary.LittleEndian.Uint32(raw[position : position+4]),
		Flags:         binary.LittleEndian.Uint16(raw[position+4 : position+6]),
		Version:       binary.LittleEndian.Uint16(raw[position+6 : position+8]),
		Offset:        int32(binary.LittleEndian.Uint32(raw[position+8 : position+12])),
		LogicalLength: int32(binary.LittleEndian.Uint32(raw[position+12 : position+16])),
	}
}

func isGameSentinel(entry lumps.GameLump) bool {
	return entry.ID == 0 && entry.Flags == 0 && entry.Version == 0 && entry.LogicalLength == 0
}

func gameEntryRange(raw []byte, base int32, count int, directoryEnd int64, index int, entry lumps.GameLump) (int64, int64, error) {
	if entry.Offset < 0 || entry.LogicalLength < 0 {
		return 0, 0, fmt.Errorf("negative offset or logical length")
	}
	if entry.Flags&1 == 0 && entry.LogicalLength == 0 && entry.Offset == 0 {
		return 0, 0, nil
	}
	start := int64(entry.Offset) - int64(base)
	if start < 0 {
		return 0, 0, fmt.Errorf("offset %d precedes container base %d", entry.Offset, base)
	}
	if start > int64(len(raw)) {
		return 0, 0, fmt.Errorf("offset outside container: %w", io.ErrUnexpectedEOF)
	}
	end := start + int64(entry.LogicalLength)
	if entry.Flags&1 != 0 {
		if index+1 >= count {
			return 0, 0, fmt.Errorf("compressed payload lacks following boundary")
		}
		next := readGameEntry(raw, index+1)
		if index+1 == count-1 && isGameSentinel(next) && next.Offset == 0 {
			end = int64(len(raw))
		} else {
			if next.Offset <= 0 {
				return 0, 0, fmt.Errorf("compressed boundary entry %d has invalid offset %d", index+1, next.Offset)
			}
			end = int64(next.Offset) - int64(base)
		}
		if end < start {
			return 0, 0, fmt.Errorf("compressed boundary entry %d precedes payload", index+1)
		}
	}
	if end > int64(len(raw)) {
		return 0, 0, fmt.Errorf("stored span outside container: %w", io.ErrUnexpectedEOF)
	}
	if (end > start || entry.LogicalLength != 0) && start < directoryEnd {
		return 0, 0, fmt.Errorf("payload overlaps directory")
	}
	return start, end, nil
}

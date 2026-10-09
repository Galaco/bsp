package bsp

import (
	"encoding/binary"
	"fmt"
	"github.com/galaco/bsp/internal/codec"
	"github.com/galaco/bsp/lumps"
)

func (bsp *File) decodedLumpView(index LumpID) ([]byte, error) {
	raw := bsp.lumps[index].raw
	expected := binary.LittleEndian.Uint32(bsp.header.Lumps[index].FourCC[:])
	if expected == 0 {
		if uint64(len(raw)) > bsp.limits.MaxDecodedLumpBytes {
			return nil, fmt.Errorf("%w: lump %d bytes %d", ErrLimitExceeded, index, len(raw))
		}
		return raw, nil
	}
	decoded, err := codec.DecodeLZMA(raw, expected, false, bsp.limits)
	if err != nil {
		return nil, &FormatError{fmt.Sprintf("decompress lump %d", index), err}
	}
	return decoded, nil
}

// DecodedLump returns owned logical bytes; RawLump returns exact stored bytes.
// Only explicit nonzero descriptor metadata selects Source LZMA compression.
func (bsp *File) DecodedLump(index LumpID) ([]byte, error) {
	if err := validateID(index); err != nil {
		return nil, err
	}
	decoded, err := bsp.decodedLumpView(index)
	if err != nil {
		return nil, err
	}
	if binary.LittleEndian.Uint32(bsp.header.Lumps[index].FourCC[:]) == 0 {
		return cloneBytes(decoded), nil
	}
	return decoded, nil
}

func (bsp *File) decodedGameEntry(entry lumps.GameLump) ([]byte, error) {
	if entry.Flags&1 == 0 {
		if uint64(len(entry.StoredData)) > bsp.limits.MaxDecodedLumpBytes {
			return nil, fmt.Errorf("%w: game ID %#x bytes %d", ErrLimitExceeded, entry.ID, len(entry.StoredData))
		}
		return entry.StoredData, nil
	}
	decoded, err := codec.DecodeLZMA(entry.StoredData, uint32(entry.LogicalLength), true, bsp.limits)
	if err != nil {
		return nil, &FormatError{fmt.Sprintf("decompress game ID %#x version %d", entry.ID, entry.Version), err}
	}
	return decoded, nil
}

// DecodedGameLump returns an owned logical entry snapshot without editing storage.
// Compression is selected by the entry flag, never by sniffing unmarked prefixes.
func (bsp *File) DecodedGameLump(index int) ([]byte, error) {
	directory, err := bsp.gameDirectory()
	if err != nil {
		return nil, err
	}
	if index < 0 || index >= len(directory.Entries) {
		return nil, fmt.Errorf("%w: %d", ErrInvalidGameLumpIndex, index)
	}
	entry := directory.Entries[index]
	decoded, err := bsp.decodedGameEntry(entry)
	if err != nil {
		return nil, err
	}
	if entry.Flags&1 == 0 {
		return cloneBytes(decoded), nil
	}
	return decoded, nil
}

// SetDecodedGameLump commits logical bytes uncompressed, preserving ID/version
// and every flag except the compression bit. Counts/positions/lengths are derived.
func (bsp *File) SetDecodedGameLump(index int, decoded []byte) error {
	directory, err := bsp.gameDirectory()
	if err != nil {
		return err
	}
	if index < 0 || index >= len(directory.Entries) {
		return fmt.Errorf("%w: %d", ErrInvalidGameLumpIndex, index)
	}
	return bsp.commitDecodedGameEntry(directory, index, decoded)
}

func (bsp *File) commitDecodedGameEntry(directory codec.GameDirectory, index int, decoded []byte) error {
	entry := &directory.Entries[index]
	var tail []byte
	if entry.Flags&1 != 0 {
		end, err := codec.LZMAStoredSize(entry.StoredData, uint32(entry.LogicalLength), true)
		if err != nil {
			return &FormatError{"replace compressed game entry", err}
		}
		tail = entry.StoredData[end:]
	}
	entry.Flags &^= 1
	entry.StoredData = decoded
	return bsp.commitGameDirectoryEdit(directory, index, tail)
}

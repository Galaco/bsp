package bsp_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"testing"

	"github.com/galaco/bsp"
	"github.com/galaco/bsp/lumps"
)

// These fixtures write wire fields directly, independently of the game encoder.
func putGameEntry(raw []byte, index int, id uint32, flags, version uint16, offset, length int32) {
	entry := raw[4+16*index : 4+16*(index+1)]
	binary.LittleEndian.PutUint32(entry, id)
	binary.LittleEndian.PutUint16(entry[4:], flags)
	binary.LittleEndian.PutUint16(entry[6:], version)
	binary.LittleEndian.PutUint32(entry[8:], uint32(offset))
	binary.LittleEndian.PutUint32(entry[12:], uint32(length))
}

func readGameFixture(t *testing.T, raw []byte) *bsp.Bsp {
	t.Helper()
	file, err := bsp.ReadFromStream(bytes.NewReader(containerBytes(35, raw)), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	if err != nil {
		t.Fatal(err)
	}
	return file
}

func requireGameEntries(t *testing.T, file *bsp.Bsp, count int) []lumps.GameLump {
	t.Helper()
	entries, err := file.GameLumps()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != count {
		t.Fatalf("game entry count=%d, want %d", len(entries), count)
	}
	return entries
}

func requireGameEntry(t *testing.T, got, want lumps.GameLump) {
	t.Helper()
	if got.ID != want.ID || got.Flags != want.Flags || got.Version != want.Version || got.Offset != want.Offset || got.LogicalLength != want.LogicalLength || !bytes.Equal(got.StoredData, want.StoredData) {
		t.Fatalf("game entry=%+v, want %+v", got, want)
	}
}

func requireStoredGame(t *testing.T, file *bsp.Bsp, want []byte) {
	t.Helper()
	for _, getter := range []struct {
		name string
		get  func() ([]byte, error)
	}{
		{"raw", func() ([]byte, error) { return file.RawLump(35) }},
		{"writer", func() ([]byte, error) { return bsp.NewWriter(file).WriteLump(35) }},
	} {
		got, err := getter.get()
		if err != nil || !bytes.Equal(got, want) {
			t.Fatalf("%s game bytes=%x, want %x; error=%v", getter.name, got, want, err)
		}
	}
}

func requireGameRelocation(t *testing.T, file *bsp.Bsp, original []byte, base int32) []byte {
	t.Helper()
	got, err := file.MarshalGameLump(base)
	if err != nil {
		t.Fatal(err)
	}
	want := append([]byte(nil), original...)
	if len(want) != 0 {
		count := int(binary.LittleEndian.Uint32(want))
		for i := 0; i < count; i++ {
			field := want[12+16*i : 16+16*i]
			offset := int32(binary.LittleEndian.Uint32(field))
			if offset != 0 {
				binary.LittleEndian.PutUint32(field, uint32(offset-bsp.HeaderSize+base))
			}
		}
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("relocation to %d changed more than offsets: got %x, want %x", base, got, want)
	}
	return got
}

func TestGameLumpMinimalWireHasNoPayloadLengthPrefix(t *testing.T) {
	// Count=1, ID=0xdeadbeef, flags=0xa4, version=0x1234,
	// absolute offset=1036+20, length=3; the next byte is payload, not a length.
	raw := []byte{
		1, 0, 0, 0,
		0xef, 0xbe, 0xad, 0xde, 0xa4, 0, 0x34, 0x12,
		0x20, 0x04, 0, 0, 3, 0, 0, 0,
		0xaa, 0xbb, 0xcc,
	}
	file := readGameFixture(t, raw)
	want := lumps.GameLump{ID: 0xdeadbeef, Flags: 0xa4, Version: 0x1234, Offset: 1056, LogicalLength: 3, StoredData: []byte{0xaa, 0xbb, 0xcc}}
	requireGameEntry(t, requireGameEntries(t, file, 1)[0], want)
	header := file.Header()
	for _, base := range []int32{2048, 4096, 2048} {
		emitted := requireGameRelocation(t, file, raw, base)
		if len(emitted) != 23 || !bytes.Equal(emitted[20:], want.StoredData) {
			t.Fatalf("minimal game wire=%x (length %d), want 23 bytes with payload at 20", emitted, len(emitted))
		}
		// Embed the relocated bytes at their actual absolute base and consume again.
		container := make([]byte, int(base)+len(emitted))
		copy(container, containerBytes(35, nil))
		binary.LittleEndian.PutUint32(container[8+35*16:], uint32(base))
		binary.LittleEndian.PutUint32(container[12+35*16:], uint32(len(emitted)))
		copy(container[int(base):], emitted)
		reloaded, err := bsp.ReadFromStream(bytes.NewReader(container), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
		if err != nil {
			t.Fatal(err)
		}
		moved := want
		moved.Offset = base + 20
		requireGameEntry(t, requireGameEntries(t, reloaded, 1)[0], moved)
	}
	requireGameEntry(t, requireGameEntries(t, file, 1)[0], want)
	requireStoredGame(t, file, raw)
	if file.Header() != header {
		t.Fatal("repeated relocation changed source header")
	}
	// A canonical replacement must also omit any per-payload length word.
	if err := file.SetGameLumps([]lumps.GameLump{want}); err != nil {
		t.Fatal(err)
	}
	canonical := append([]byte(nil), raw...)
	binary.LittleEndian.PutUint32(canonical[12:], 20)
	requireStoredGame(t, file, canonical)
	descriptor, err := file.Descriptor(35)
	if err != nil || descriptor.Offset != 0 || descriptor.Length != 23 {
		t.Fatalf("replacement descriptor=%+v, error=%v", descriptor, err)
	}
	emitted, err := file.MarshalGameLump(2048)
	if err != nil || len(emitted) != 23 || binary.LittleEndian.Uint32(emitted[12:]) != 2068 || !bytes.Equal(emitted[20:], want.StoredData) {
		t.Fatalf("replacement relocation=%x, error=%v", emitted, err)
	}
}

func mixedGameFixture() ([]byte, []lumps.GameLump) {
	raw := make([]byte, 81)
	binary.LittleEndian.PutUint32(raw, 4)
	copy(raw[68:], []byte("gap!abcDEwxyz!"))
	entries := []lumps.GameLump{
		{ID: 0x70727073, Flags: 0x8000, Version: 9, Offset: 1108, LogicalLength: 3, StoredData: []byte("abc")},
		{ID: 0xdecafbad, Flags: 2, Version: 513, Offset: 1111, LogicalLength: 2, StoredData: []byte("DE")},
		{ID: 0xdecafbad, Flags: 0x1000, Version: 65535, Offset: 1113, LogicalLength: 4, StoredData: []byte("wxyz")},
		{ID: 0x12345678, Flags: 4, Version: 7},
	}
	for i, entry := range entries {
		putGameEntry(raw, i, entry.ID, entry.Flags, entry.Version, entry.Offset, entry.LogicalLength)
	}
	return raw, entries
}

func TestGameLumpOrderedDuplicatesAndOwnedSnapshots(t *testing.T) {
	raw, want := mixedGameFixture()
	file := readGameFixture(t, raw)
	entries := requireGameEntries(t, file, len(want))
	for i := range want {
		requireGameEntry(t, entries[i], want[i])
		entry, err := file.GameLump(i)
		if err != nil {
			t.Fatal(err)
		}
		requireGameEntry(t, entry, want[i])
	}
	if _, err := file.GameLumpByID(0xdecafbad); !errors.Is(err, bsp.ErrAmbiguousGameLumpID) {
		t.Fatalf("duplicate ID lookup error=%v", err)
	}
	if _, err := file.GameLumpByID(0xfedcba98); !errors.Is(err, bsp.ErrLumpNotFound) {
		t.Fatalf("missing ID lookup error=%v", err)
	}
	unique, err := file.GameLumpByID(want[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	requireGameEntry(t, unique, want[0])
	entries[0].ID = 1
	entries[0].StoredData[0] = 'X'
	unique.StoredData[1] = 'Y'
	indexed, err := file.GameLump(2)
	if err != nil {
		t.Fatal(err)
	}
	indexed.StoredData[0] = 'Z'
	for i, entry := range requireGameEntries(t, file, len(want)) {
		requireGameEntry(t, entry, want[i])
	}
	header := file.Header()
	for _, index := range []int{-1, len(want), math.MaxInt32} {
		if _, err := file.GameLump(index); !errors.Is(err, bsp.ErrInvalidGameLumpIndex) {
			t.Fatalf("getter index %d error=%v", index, err)
		}
		if err := file.SetGameLump(index, lumps.GameLump{ID: 1, StoredData: []byte("replacement")}); !errors.Is(err, bsp.ErrInvalidGameLumpIndex) {
			t.Fatalf("setter index %d error=%v", index, err)
		}
	}
	requireGameRelocation(t, file, raw, 2048)
	requireStoredGame(t, file, raw)
	if file.Header() != header {
		t.Fatal("failed index edits changed source metadata")
	}
}

func TestGameLumpReplacementPreservesSiblingsAndIsAtomic(t *testing.T) {
	raw, original := mixedGameFixture()
	file := readGameFixture(t, raw)
	payload := []byte("replacement")
	replacement := lumps.GameLump{ID: original[1].ID, Flags: 0x20, Version: 88, LogicalLength: 999, StoredData: payload}
	if err := file.SetGameLump(1, replacement); err != nil {
		t.Fatal(err)
	}
	payload[0] = 'X'
	entries := requireGameEntries(t, file, len(original))
	expected := append([]lumps.GameLump(nil), original...)
	expected[1] = replacement
	expected[1].StoredData = []byte("replacement")
	expected[1].LogicalLength = int32(len(expected[1].StoredData))
	for i := range expected {
		// Reconstruction may move siblings, but not change their metadata/data.
		expected[i].Offset = entries[i].Offset
		requireGameEntry(t, entries[i], expected[i])
	}
	if _, err := file.GameLumpByID(original[1].ID); !errors.Is(err, bsp.ErrAmbiguousGameLumpID) {
		t.Fatalf("replacement lost duplicate ID: %v", err)
	}
	if err := file.SetGameLumps(entries); err != nil {
		t.Fatal(err)
	}
	entries[0].ID = 0
	entries[0].StoredData[0] = 'Y'
	current := requireGameEntries(t, file, len(expected))
	for i := range current {
		expected[i].Offset = current[i].Offset
		requireGameEntry(t, current[i], expected[i])
	}
	before, err := file.RawLump(35)
	if err != nil {
		t.Fatal(err)
	}
	descriptor, err := file.Descriptor(35)
	if err != nil || descriptor.Offset != 0 || descriptor.Length != int32(len(before)) {
		t.Fatalf("canonical descriptor=%+v, error=%v", descriptor, err)
	}
	// A last compressed entry has no stored-span boundary and cannot be committed.
	invalid := lumps.GameLump{ID: 0x55555555, Flags: 1, Version: 91, LogicalLength: 100, StoredData: []byte("opaque")}
	for _, edit := range []struct {
		name string
		set  func() error
	}{
		{"entry", func() error { return file.SetGameLump(len(current)-1, invalid) }},
		{"directory", func() error { return file.SetGameLumps([]lumps.GameLump{invalid}) }},
	} {
		if err := edit.set(); !errors.Is(err, bsp.ErrMalformed) {
			t.Fatalf("invalid %s edit error=%v", edit.name, err)
		}
		requireStoredGame(t, file, before)
		gotDescriptor, err := file.Descriptor(35)
		if err != nil || gotDescriptor != descriptor {
			t.Fatalf("failed %s edit changed metadata: %+v, error=%v", edit.name, gotDescriptor, err)
		}
		for i, entry := range requireGameEntries(t, file, len(expected)) {
			requireGameEntry(t, entry, expected[i])
		}
	}
}

func TestGameLumpRelocationPreservesSharedViewsAndUnclaimedBytes(t *testing.T) {
	raw := make([]byte, 104)
	for i := range raw {
		raw[i] = byte(0x80 + i)
	}
	binary.LittleEndian.PutUint32(raw, 5)
	putGameEntry(raw, 0, 0x11111111, 0x8000, 17, 1124, 5)
	putGameEntry(raw, 1, 0x22222222, 2, 18, 1126, 3)
	putGameEntry(raw, 2, 0x33333333, 4, 19, 1124, 5)
	putGameEntry(raw, 3, 0x44444444, 8, 20, 0, 0)
	putGameEntry(raw, 4, 0x55555555, 16, 21, 1140, 0)
	file := readGameFixture(t, raw)
	entries := requireGameEntries(t, file, 5)
	if !bytes.Equal(entries[0].StoredData, raw[88:93]) || !bytes.Equal(entries[1].StoredData, raw[90:93]) || !bytes.Equal(entries[2].StoredData, raw[88:93]) || len(entries[3].StoredData) != 0 || len(entries[4].StoredData) != 0 {
		t.Fatalf("shared and empty views=%+v", entries)
	}
	entries[0].StoredData[2] ^= 0xff
	if !bytes.Equal(entries[1].StoredData, raw[90:93]) || !bytes.Equal(entries[2].StoredData, raw[88:93]) {
		t.Fatal("separate public snapshots alias overlapping views")
	}
	header := file.Header()
	for _, base := range []int32{2048, 4096, 0, 2048} {
		emitted := requireGameRelocation(t, file, raw, base)
		if !bytes.Equal(emitted[84:], raw[84:]) || binary.LittleEndian.Uint32(emitted[60:]) != 0 || binary.LittleEndian.Uint32(emitted[76:]) != uint32(base+104) {
			t.Fatalf("gaps, trailing bytes or empty offsets changed: %x", emitted)
		}
		emitted[len(emitted)-1] ^= 0xff
	}
	for _, base := range []int32{-1, math.MaxInt32} {
		if _, err := file.MarshalGameLump(base); !errors.Is(err, bsp.ErrMalformed) {
			t.Fatalf("invalid relocation base %d error=%v", base, err)
		}
	}
	unchanged := requireGameEntries(t, file, 5)
	if !bytes.Equal(unchanged[0].StoredData, raw[88:93]) || unchanged[4].Offset != 1140 {
		t.Fatalf("source views changed after export: %+v", unchanged)
	}
	requireStoredGame(t, file, raw)
	if file.Header() != header {
		t.Fatal("relocation changed imported header")
	}
}

func TestGameLumpAbsentAndCountZeroRemainDistinct(t *testing.T) {
	for _, tc := range []struct {
		name string
		id   int
		raw  []byte
	}{
		{"absent", -1, nil},
		{"raw empty", 35, []byte{}},
		{"present zero count", 35, []byte{0, 0, 0, 0}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			file, err := bsp.ReadFromStream(bytes.NewReader(containerBytes(tc.id, tc.raw)), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
			if err != nil {
				t.Fatal(err)
			}
			requireGameEntries(t, file, 0)
			if _, err := file.GameLumpByID(0); !errors.Is(err, bsp.ErrLumpNotFound) {
				t.Fatalf("empty ID lookup error=%v", err)
			}
			if _, err := file.GameLump(0); !errors.Is(err, bsp.ErrInvalidGameLumpIndex) {
				t.Fatalf("empty index lookup error=%v", err)
			}
			requireGameRelocation(t, file, tc.raw, 2048)
			requireStoredGame(t, file, tc.raw)
			descriptor, err := file.Descriptor(35)
			if err != nil || descriptor.Length != int32(len(tc.raw)) {
				t.Fatalf("empty descriptor=%+v, error=%v", descriptor, err)
			}
		})
	}
}

func TestGameLumpCompressedStoredSpansAndTerminalConventions(t *testing.T) {
	for _, explicit := range []bool{false, true} {
		name := "zero terminal"
		if explicit {
			name = "explicit terminal"
		}
		t.Run(name, func(t *testing.T) {
			raw := make([]byte, 70)
			binary.LittleEndian.PutUint32(raw, 3)
			copy(raw[52:], []byte("gap!ABCDE12345TAIL"))
			putGameEntry(raw, 0, 0xabcdef01, 0x8001, 42, 1092, 100000)
			putGameEntry(raw, 1, 0xabcdef02, 1, 43, 1097, 500)
			terminal := int32(0)
			end := len(raw)
			if explicit {
				terminal = 1102
				end = 66
			}
			putGameEntry(raw, 2, 0, 0, 0, terminal, 0)
			file := readGameFixture(t, raw)
			want := []lumps.GameLump{
				{ID: 0xabcdef01, Flags: 0x8001, Version: 42, Offset: 1092, LogicalLength: 100000, StoredData: raw[56:61]},
				{ID: 0xabcdef02, Flags: 1, Version: 43, Offset: 1097, LogicalLength: 500, StoredData: raw[61:end]},
				{Offset: terminal},
			}
			for i, entry := range requireGameEntries(t, file, 3) {
				requireGameEntry(t, entry, want[i])
			}
			if _, err := file.GameLumpByID(0); !errors.Is(err, bsp.ErrLumpNotFound) {
				t.Fatalf("terminal sentinel participated in ID lookup: %v", err)
			}
			for _, base := range []int32{2048, 4096} {
				requireGameRelocation(t, file, raw, base)
			}
			for i, entry := range requireGameEntries(t, file, 3) {
				requireGameEntry(t, entry, want[i])
			}
			requireStoredGame(t, file, raw)
			// Arbitrary stored bytes are opaque here, not a decompression request.
			if err := file.SetGameLumps(want); err != nil {
				t.Fatal(err)
			}
			canonical := make([]byte, 52+len(want[0].StoredData)+len(want[1].StoredData))
			binary.LittleEndian.PutUint32(canonical, 3)
			putGameEntry(canonical, 0, want[0].ID, want[0].Flags, want[0].Version, 52, 100000)
			putGameEntry(canonical, 1, want[1].ID, want[1].Flags, want[1].Version, 57, 500)
			canonicalEnd := int32(0)
			if explicit {
				canonicalEnd = int32(len(canonical))
			}
			putGameEntry(canonical, 2, 0, 0, 0, canonicalEnd, 0)
			copy(canonical[52:], want[0].StoredData)
			copy(canonical[57:], want[1].StoredData)
			requireStoredGame(t, file, canonical)
			current := requireGameEntries(t, file, 3)
			for i := range want {
				want[i].Offset = current[i].Offset
				requireGameEntry(t, current[i], want[i])
			}
		})
	}
}

func TestGameLumpMalformedDirectoriesAreContextualAndNonMutating(t *testing.T) {
	one := func(offset, length int32, flags uint16) []byte {
		raw := make([]byte, 23)
		binary.LittleEndian.PutUint32(raw, 1)
		putGameEntry(raw, 0, 0xdeadbeef, flags, 3, offset, length)
		copy(raw[20:], "abc")
		return raw
	}
	compressedBoundary := func(offset int32, sentinel bool) []byte {
		raw := make([]byte, 42)
		binary.LittleEndian.PutUint32(raw, 2)
		putGameEntry(raw, 0, 0xdeadbeef, 1, 3, 1072, 1000)
		id := uint32(0)
		if !sentinel {
			id = 0x12345678
		}
		putGameEntry(raw, 1, id, 0, 0, offset, 0)
		return raw
	}
	truncatedDirectory := make([]byte, 19)
	binary.LittleEndian.PutUint32(truncatedDirectory, 1)
	for _, tc := range []struct {
		name string
		raw  []byte
		eof  bool
	}{
		{"one byte header", []byte{1}, true},
		{"two byte header", []byte{1, 0}, true},
		{"three byte header", []byte{1, 0, 0}, true},
		{"negative count", []byte{0xff, 0xff, 0xff, 0xff}, false},
		{"overflowing directory", []byte{0xff, 0xff, 0xff, 0x7f}, false},
		{"truncated directory", truncatedDirectory, true},
		{"negative offset", one(-1, 3, 0), false},
		{"negative logical length", one(1056, -1, 0), false},
		{"before container base", one(1035, 3, 0), false},
		{"payload overlaps directory", one(1055, 3, 0), false},
		{"payload starts past end", one(1060, 1, 0), false},
		{"payload extends past end", one(1056, 4, 0), false},
		{"huge logical span", one(1056, math.MaxInt32, 0), false},
		{"nonempty zero offset", one(0, 3, 0), false},
		{"empty offset past end", one(1060, 0, 0), false},
		{"compressed missing boundary", one(1056, 1000, 1), false},
		{"compressed reversed boundary", compressedBoundary(1071, true), false},
		{"compressed boundary past end", compressedBoundary(1079, true), false},
		{"compressed negative boundary", compressedBoundary(-1, true), false},
		{"compressed zero nonsentinel boundary", compressedBoundary(0, false), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			file := readGameFixture(t, tc.raw)
			header := file.Header()
			for _, operation := range []struct {
				name string
				call func() error
			}{
				{"list", func() error { _, err := file.GameLumps(); return err }},
				{"index", func() error { _, err := file.GameLump(0); return err }},
				{"ID", func() error { _, err := file.GameLumpByID(0xdeadbeef); return err }},
				{"relocate", func() error { _, err := file.MarshalGameLump(2048); return err }},
				{"replace entry", func() error { return file.SetGameLump(0, lumps.GameLump{ID: 1, StoredData: []byte("valid")}) }},
			} {
				err := operation.call()
				var formatError *bsp.FormatError
				if !errors.Is(err, bsp.ErrMalformed) || !errors.As(err, &formatError) {
					t.Fatalf("%s malformed error lacks classification/context: %v", operation.name, err)
				}
				if tc.eof && !errors.Is(err, io.ErrUnexpectedEOF) {
					t.Fatalf("%s truncation error=%v, want unexpected EOF", operation.name, err)
				}
				requireStoredGame(t, file, tc.raw)
				if file.Header() != header {
					t.Fatalf("%s changed malformed source metadata", operation.name)
				}
			}
		})
	}
}

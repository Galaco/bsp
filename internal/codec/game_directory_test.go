package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/galaco/bsp/lumps"
)

// wireGameEntry and authoredGameWire deliberately do not use the codec encoder.
// Entries occupy exactly 16 bytes, with signed fields authored as their wire bits.
type wireGameEntry struct {
	id      uint32
	flags   uint16
	version uint16
	offset  int32
	length  int32
}

func authoredGameWire(entries []wireGameEntry, tail []byte) []byte {
	raw := make([]byte, 4+16*len(entries)+len(tail))
	binary.LittleEndian.PutUint32(raw[:4], uint32(len(entries)))
	for index, entry := range entries {
		position := 4 + index*16
		binary.LittleEndian.PutUint32(raw[position:position+4], entry.id)
		binary.LittleEndian.PutUint16(raw[position+4:position+6], entry.flags)
		binary.LittleEndian.PutUint16(raw[position+6:position+8], entry.version)
		binary.LittleEndian.PutUint32(raw[position+8:position+12], uint32(entry.offset))
		binary.LittleEndian.PutUint32(raw[position+12:position+16], uint32(entry.length))
	}
	copy(raw[4+16*len(entries):], tail)
	return raw
}

func TestGameDirectoryAuthoredViewsAndRelocation(t *testing.T) {
	const base int32 = 1000
	wireEntries := []wireGameEntry{
		{0xfedcba98, 0x8002, 37, base + 88, 6},
		{0xfedcba98, 0, 9, base + 90, 4},
		{0x01020304, 0x4001, 22, base + 94, 4000},
		{0x99887766, 0, 10, base + 102, 3},
		{0, 0, 0, base + 110, 0},
	}
	raw := authoredGameWire(wireEntries, []byte("gap!ABCDEFpacked!!XYZtrailer"))
	original := append([]byte(nil), raw...)
	directory, err := ParseGameDirectory(raw, base)
	if err != nil {
		t.Fatal(err)
	}
	wantData := []string{"ABCDEF", "CDEF", "packed!!", "XYZ", ""}
	for index, entry := range directory.Entries {
		wire := wireEntries[index]
		if entry.ID != wire.id || entry.Flags != wire.flags || entry.Version != wire.version || entry.Offset != wire.offset || entry.LogicalLength != wire.length {
			t.Fatalf("entry %d metadata = %#v, wire = %#v", index, entry, wire)
		}
		if string(entry.StoredData) != wantData[index] {
			t.Fatalf("entry %d data = %q", index, entry.StoredData)
		}
	}
	if &directory.Raw[0] != &raw[0] || &directory.Entries[0].StoredData[2] != &directory.Entries[1].StoredData[0] {
		t.Fatal("parser did not borrow shared raw views")
	}
	for _, newBase := range []int32{2048, 8192} {
		for repeat := 0; repeat < 2; repeat++ {
			moved, err := RelocateGameDirectory(directory, newBase)
			if err != nil {
				t.Fatal(err)
			}
			wantEntries := append([]wireGameEntry(nil), wireEntries...)
			for index := range wantEntries {
				wantEntries[index].offset += newBase - base
			}
			want := authoredGameWire(wantEntries, original[84:])
			if !bytes.Equal(moved, want) {
				t.Fatalf("relocation to %d changed non-offset bytes", newBase)
			}
			parsed, err := ParseGameDirectory(moved, newBase)
			if err != nil {
				t.Fatal(err)
			}
			back, err := RelocateGameDirectory(parsed, base)
			if err != nil || !bytes.Equal(back, original) {
				t.Fatalf("reparse/return: %v, bytes equal = %v", err, bytes.Equal(back, original))
			}
			if !bytes.Equal(raw, original) {
				t.Fatal("relocation mutated input")
			}
		}
	}
}

func TestGameDirectoryEmptyFormsAndOffsets(t *testing.T) {
	const base int32 = 64
	for _, raw := range [][]byte{nil, {}, {0, 0, 0, 0}, {0, 0, 0, 0, 0xaa, 0xbb}} {
		directory, err := ParseGameDirectory(raw, base)
		if err != nil || len(directory.Entries) != 0 || !bytes.Equal(directory.Raw, raw) {
			t.Fatalf("empty parse: %#v, %v", directory, err)
		}
		moved, err := RelocateGameDirectory(directory, 128)
		if err != nil || !bytes.Equal(moved, raw) || len(moved) != len(raw) {
			t.Fatalf("empty relocation: %x, %v", moved, err)
		}
	}
	encoded, err := EncodeGameDirectory(nil, base)
	if err != nil || !bytes.Equal(encoded, []byte{0, 0, 0, 0}) {
		t.Fatalf("canonical empty = %x, %v", encoded, err)
	}
	entries := []wireGameEntry{{1, 0, 2, 0, 0}, {2, 0, 3, base, 0}, {3, 0, 4, base + 52, 0}}
	raw := authoredGameWire(entries, nil)
	directory, err := ParseGameDirectory(raw, base)
	if err != nil {
		t.Fatal(err)
	}
	moved, err := RelocateGameDirectory(directory, 0)
	if err != nil {
		t.Fatal(err)
	}
	entries[1].offset = 0
	entries[2].offset = 52
	if !bytes.Equal(moved, authoredGameWire(entries, nil)) {
		t.Fatalf("empty offsets = %x", moved)
	}
	if _, err := ParseGameDirectory(moved, 0); err != nil {
		t.Fatal(err)
	}
}

func TestGameDirectoryUncompressedOrderingAndOverlap(t *testing.T) {
	const base int32 = 256
	raw := authoredGameWire([]wireGameEntry{
		{7, 0, 1, base + 58, 2},
		{8, 0, 2, base + 52, 8},
		{7, 0, 3, base + 52, 8},
	}, []byte("abcdefgh"))
	directory, err := ParseGameDirectory(raw, base)
	if err != nil {
		t.Fatal(err)
	}
	if string(directory.Entries[0].StoredData) != "gh" || string(directory.Entries[1].StoredData) != "abcdefgh" || &directory.Entries[1].StoredData[0] != &directory.Entries[2].StoredData[0] {
		t.Fatal("uncompressed out-of-order/shared views were not preserved")
	}
}

func TestGameDirectoryCompressedTerminalConventions(t *testing.T) {
	const base int32 = 300
	for _, boundary := range []int32{0, base + 40} {
		raw := authoredGameWire([]wireGameEntry{{11, 1, 5, base + 36, 10000}, {0, 0, 0, boundary, 0}}, []byte("zip!tail"))
		directory, err := ParseGameDirectory(raw, base)
		if err != nil {
			t.Fatal(err)
		}
		want := "zip!"
		if boundary == 0 {
			want = "zip!tail"
		}
		if directory.Entries[0].LogicalLength != 10000 || string(directory.Entries[0].StoredData) != want || len(directory.Entries[1].StoredData) != 0 {
			t.Fatalf("boundary %d: %#v", boundary, directory.Entries)
		}
		moved, err := RelocateGameDirectory(directory, 700)
		if err != nil {
			t.Fatal(err)
		}
		wantBoundary := boundary
		if boundary != 0 {
			wantBoundary += 400
		}
		if got := int32(binary.LittleEndian.Uint32(moved[28:32])); got != wantBoundary {
			t.Fatalf("terminal moved to %d, want %d", got, wantBoundary)
		}
		reparsed, err := ParseGameDirectory(moved, 700)
		if err != nil || string(reparsed.Entries[0].StoredData) != want {
			t.Fatalf("reparse: %#v, %v", reparsed, err)
		}
		encoded, err := EncodeGameDirectory(directory.Entries, 900)
		if err != nil {
			t.Fatal(err)
		}
		end := int32(0)
		if boundary != 0 {
			end = 940
		}
		wantRaw := authoredGameWire([]wireGameEntry{{11, 1, 5, 936, 10000}, {0, 0, 0, end, 0}}, []byte(want))
		if !bytes.Equal(encoded, wantRaw) {
			t.Fatalf("compressed canonical bytes = %x, want %x", encoded, wantRaw)
		}
	}
}

func TestEncodeGameDirectoryCanonicalPayloads(t *testing.T) {
	entries := []lumps.GameLump{
		{ID: 0xabcdef01, Flags: 0x8002, Version: 12, Offset: 1, LogicalLength: 999, StoredData: []byte("abc")},
		{ID: 0xabcdef01, Version: 13, StoredData: []byte("XY")},
		{ID: 9, Version: 99},
	}
	before := append([]lumps.GameLump(nil), entries...)
	for index := range before {
		before[index].StoredData = append([]byte(nil), before[index].StoredData...)
	}
	for _, base := range []int32{1024, 4096} {
		raw, err := EncodeGameDirectory(entries, base)
		if err != nil {
			t.Fatal(err)
		}
		want := authoredGameWire([]wireGameEntry{
			{0xabcdef01, 0x8002, 12, base + 52, 3},
			{0xabcdef01, 0, 13, base + 56, 2},
			{9, 0, 99, 0, 0},
		}, []byte("abc\x00XY\x00\x00"))
		if !bytes.Equal(raw, want) {
			t.Fatalf("canonical bytes = %x, want %x", raw, want)
		}
		parsed, err := ParseGameDirectory(raw, base)
		if err != nil {
			t.Fatal(err)
		}
		again, err := EncodeGameDirectory(parsed.Entries, base)
		if err != nil || !bytes.Equal(again, raw) {
			t.Fatalf("repeat encode = %x, %v", again, err)
		}
	}
	if !reflect.DeepEqual(entries, before) {
		t.Fatal("encode mutated caller metadata or payloads")
	}
	// Empty real entries still delimit their compressed predecessor.
	compressed := []lumps.GameLump{{ID: 1, Flags: 1, LogicalLength: 200, StoredData: []byte("zip")}, {ID: 2}}
	raw, err := EncodeGameDirectory(compressed, 2000)
	want := authoredGameWire([]wireGameEntry{{1, 1, 0, 2036, 200}, {2, 0, 0, 2040, 0}}, []byte("zip\x00"))
	if err != nil || !bytes.Equal(raw, want) {
		t.Fatalf("real empty boundary: %x, %v", raw, err)
	}
	if _, err := ParseGameDirectory(raw, 2000); err != nil {
		t.Fatal(err)
	}
}

func TestParseGameDirectoryMalformed(t *testing.T) {
	const base int32 = 100
	cases := []struct {
		name string
		raw  []byte
		base int32
		eof  bool
	}{
		{"one-byte", []byte{0}, base, true},
		{"two-byte", []byte{0, 0}, base, true},
		{"three-byte", []byte{0, 0, 0}, base, true},
		{"negative-count", []byte{255, 255, 255, 255}, base, false},
		{"overflowing-count", []byte{255, 255, 255, 127}, base, false},
		{"truncated-directory", []byte{1, 0, 0, 0}, base, true},
		{"negative-base", nil, -1, false},
		{"extent-overflow", []byte{0, 0, 0, 0}, 1<<31 - 2, false},
		{"negative-offset", authoredGameWire([]wireGameEntry{{1, 0, 0, -1, 1}}, []byte{1}), base, false},
		{"negative-length", authoredGameWire([]wireGameEntry{{1, 0, 0, base + 20, -1}}, []byte{1}), base, false},
		{"before-base", authoredGameWire([]wireGameEntry{{1, 0, 0, base - 1, 1}}, []byte{1}), base, false},
		{"directory-overlap", authoredGameWire([]wireGameEntry{{1, 0, 0, base + 19, 1}}, []byte{1}), base, false},
		{"past-end", authoredGameWire([]wireGameEntry{{1, 0, 0, base + 21, 1}}, []byte{1}), base, true},
		{"wide-end", authoredGameWire([]wireGameEntry{{1, 0, 0, base + 20, 1<<31 - 1}}, []byte{1}), base, true},
		{"empty-before-base", authoredGameWire([]wireGameEntry{{1, 0, 0, base - 1, 0}}, nil), base, false},
		{"empty-past-end", authoredGameWire([]wireGameEntry{{1, 0, 0, base + 21, 0}}, nil), base, true},
		{"missing-boundary", authoredGameWire([]wireGameEntry{{1, 1, 0, base + 20, 1000}}, []byte{1}), base, false},
		{"reversed-boundary", authoredGameWire([]wireGameEntry{{1, 1, 0, base + 38, 1000}, {2, 0, 0, base + 36, 1}}, []byte{1, 2, 3}), base, false},
		{"zero-real-boundary", authoredGameWire([]wireGameEntry{{1, 1, 0, base + 36, 1000}, {2, 0, 0, 0, 0}}, []byte{1}), base, false},
		{"negative-boundary", authoredGameWire([]wireGameEntry{{1, 1, 0, base + 36, 1000}, {0, 0, 0, -1, 0}}, []byte{1}), base, false},
		{"outside-boundary", authoredGameWire([]wireGameEntry{{1, 1, 0, base + 36, 1000}, {0, 0, 0, base + 38, 0}}, []byte{1}), base, true},
		{"compressed-directory-overlap", authoredGameWire([]wireGameEntry{{1, 1, 0, base + 35, 1000}, {0, 0, 0, 0, 0}}, []byte{1}), base, false},
		{"nonterminal-zero-entry", authoredGameWire([]wireGameEntry{{1, 1, 0, base + 52, 1000}, {0, 0, 0, 0, 0}, {2, 0, 0, 0, 0}}, []byte{1}), base, false},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, err := ParseGameDirectory(test.raw, test.base)
			if err == nil {
				t.Fatal("accepted malformed wire")
			}
			if test.eof && !errors.Is(err, io.ErrUnexpectedEOF) {
				t.Fatalf("missing EOF cause: %v", err)
			}
		})
	}
	// Every incomplete prefix of a valid nonempty one-entry payload is invalid,
	// except the deliberately valid absent form.
	raw := authoredGameWire([]wireGameEntry{{3, 0, 1, base + 20, 3}}, []byte("abc"))
	for length := 1; length < len(raw); length++ {
		if _, err := ParseGameDirectory(raw[:length], base); err == nil {
			t.Fatalf("accepted truncated prefix of length %d", length)
		}
	}
}

func TestGameDirectoryEncodeAndRelocateFailureAtomicity(t *testing.T) {
	for _, entries := range [][]lumps.GameLump{
		{{ID: 1, Offset: -1, StoredData: []byte("abc")}},
		{{ID: 1, Flags: 1, LogicalLength: -1}, {}},
		{{ID: 1, Flags: 1, LogicalLength: 3, StoredData: []byte("abc")}},
		{{StoredData: []byte("not a sentinel payload")}},
	} {
		before := append([]lumps.GameLump(nil), entries...)
		for index := range before {
			before[index].StoredData = append([]byte(nil), before[index].StoredData...)
		}
		if raw, err := EncodeGameDirectory(entries, 0); err == nil || raw != nil {
			t.Fatalf("invalid encode returned %x, %v", raw, err)
		}
		if !reflect.DeepEqual(entries, before) {
			t.Fatal("failed encode mutated input")
		}
	}
	entries := []lumps.GameLump{{ID: 1, StoredData: []byte("abc")}}
	for _, base := range []int32{-1, 1<<31 - 21} {
		if raw, err := EncodeGameDirectory(entries, base); err == nil || raw != nil {
			t.Fatalf("invalid target encode = %x, %v", raw, err)
		}
	}
	raw := authoredGameWire([]wireGameEntry{{1, 0, 0, 120, 3}}, []byte("abc"))
	before := append([]byte(nil), raw...)
	directory, err := ParseGameDirectory(raw, 100)
	if err != nil {
		t.Fatal(err)
	}
	for _, base := range []int32{-1, 1<<31 - 21} {
		if result, err := RelocateGameDirectory(directory, base); err == nil || result != nil {
			t.Fatalf("invalid relocation = %x, %v", result, err)
		}
		if !bytes.Equal(raw, before) {
			t.Fatal("failed relocation mutated source")
		}
	}
}

func TestGameDirectoryMinimalPayloadAndRawOnlyRelocation(t *testing.T) {
	entry := lumps.GameLump{ID: 0x12345678, Flags: 2, Version: 31, StoredData: []byte("abc")}
	for _, base := range []int32{1036, 8192} {
		raw, err := EncodeGameDirectory([]lumps.GameLump{entry}, base)
		if err != nil {
			t.Fatal(err)
		}
		want := authoredGameWire([]wireGameEntry{{0x12345678, 2, 31, base + 20, 3}}, []byte("abc"))
		if len(raw) != 23 || !bytes.Equal(raw, want) {
			t.Fatalf("minimal wire = %x, want %x", raw, want)
		}
		moved, err := RelocateGameDirectory(GameDirectory{Raw: raw, Base: base}, 2048)
		wantMoved := authoredGameWire([]wireGameEntry{{0x12345678, 2, 31, 2068, 3}}, []byte("abc"))
		if err != nil || !bytes.Equal(moved, wantMoved) {
			t.Fatalf("raw-only relocation = %x, %v", moved, err)
		}
	}
	malformed := GameDirectory{Raw: []byte{1, 0, 0, 0}, Base: 100}
	if raw, err := RelocateGameDirectory(malformed, 200); raw != nil || !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("raw-only malformed relocation = %x, %v", raw, err)
	}
}

func FuzzParseGameDirectory(f *testing.F) {
	f.Add([]byte{}, int32(0))
	f.Add([]byte{0, 0, 0, 0}, int32(100))
	f.Add(authoredGameWire([]wireGameEntry{{1, 0, 7, 276, 3}}, []byte("abc")), int32(256))
	f.Add(authoredGameWire([]wireGameEntry{{1, 1, 7, 292, 1000}, {0, 0, 0, 0, 0}}, []byte("zip")), int32(256))
	f.Fuzz(func(t *testing.T, raw []byte, base int32) {
		if len(raw) > 4096 {
			return
		}
		before := append([]byte(nil), raw...)
		directory, err := ParseGameDirectory(raw, base)
		if err != nil {
			return
		}
		exported, err := RelocateGameDirectory(directory, base)
		if err != nil || !bytes.Equal(exported, raw) {
			t.Fatalf("identity relocation = %x, %v", exported, err)
		}
		if !bytes.Equal(raw, before) {
			t.Fatal("parse/export mutated source")
		}
		reparsed, err := ParseGameDirectory(exported, base)
		if err != nil || !reflect.DeepEqual(directory.Entries, reparsed.Entries) {
			t.Fatalf("repeat parse: %v", err)
		}
	})
}

type gameLimitedWriter struct {
	bytes.Buffer
	remaining int
	err       error
}

func (writer *gameLimitedWriter) Write(raw []byte) (int, error) {
	if len(raw) > writer.remaining {
		n, _ := writer.Buffer.Write(raw[:writer.remaining])
		writer.remaining = 0
		return n, writer.err
	}
	n, err := writer.Buffer.Write(raw)
	writer.remaining -= n
	return n, err
}

func TestGameRelocationWriteTo(t *testing.T) {
	const base int32 = 100
	raw := authoredGameWire([]wireGameEntry{
		{1, 0, 0, base + 52, 4},
		{2, 0, 0, 0, 0},
		{3, 0, 0, base + 56, 3},
	}, []byte("ABCDxyz"))
	before := append([]byte(nil), raw...)
	directory, err := ParseGameDirectory(raw, base)
	if err != nil {
		t.Fatal(err)
	}
	writeErr := errors.New("writer failure")
	for _, newBase := range []int32{base, 800} {
		plan, err := PlanGameRelocation(directory, newBase)
		if err != nil {
			t.Fatal(err)
		}
		wantEntries := []wireGameEntry{
			{1, 0, 0, newBase + 52, 4},
			{2, 0, 0, 0, 0},
			{3, 0, 0, newBase + 56, 3},
		}
		want := authoredGameWire(wantEntries, []byte("ABCDxyz"))
		var output bytes.Buffer
		n, err := plan.WriteTo(&output)
		if err != nil || n != int64(len(want)) || !bytes.Equal(output.Bytes(), want) {
			t.Fatalf("full stream at %d: count=%d bytes=%x error=%v", newBase, n, output.Bytes(), err)
		}
		// Cuts cover ordinary spans, patched offset words and the final payload.
		for _, limit := range []int{0, 11, 12, 13, 15, 16, 43, 45, 48, 55, 58} {
			for _, injected := range []error{nil, writeErr} {
				writer := &gameLimitedWriter{remaining: limit, err: injected}
				n, err := plan.WriteTo(writer)
				expectedErr := injected
				if expectedErr == nil {
					expectedErr = io.ErrShortWrite
				}
				if n != int64(limit) || !errors.Is(err, expectedErr) || !bytes.Equal(writer.Bytes(), want[:limit]) {
					t.Fatalf("partial stream at %d/%d: count=%d bytes=%x error=%v", newBase, limit, n, writer.Bytes(), err)
				}
			}
		}
	}
	if !bytes.Equal(raw, before) {
		t.Fatal("streaming relocation mutated input")
	}
}

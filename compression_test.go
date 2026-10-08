package bsp_test

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/galaco/bsp"
	"github.com/galaco/bsp/lumps"
	"testing"
)

// Independently generated with CPython/liblzma raw LZMA1 (EOS), dictionary4096,
// lc3/lp0/pb2, then packed into Valve's17-byte envelope. Plaintext is two
// independently authored20-byte planes: normals X/Y, distances12/-9, axes0/1.
func compressedPlaneFixture(t *testing.T) []byte {
	t.Helper()
	raw, err := hex.DecodeString("4c5a4d41280000001b0000005d00100000000061fef5fedd9606c99ba4d003ff2520a2d803b3ffff63a40000")
	if err != nil {
		t.Fatal(err)
	}
	return raw
}
func readCompressedPlanes(t *testing.T, options bsp.ReadOptions) *bsp.File {
	t.Helper()
	raw := containerBytes(1, compressedPlaneFixture(t))
	binary.LittleEndian.PutUint32(raw[8+16+12:], 40)
	file, err := bsp.Read(bytes.NewReader(raw), options)
	if err != nil {
		t.Fatal(err)
	}
	return file
}
func requireTwoPlanes(t *testing.T, file *bsp.File, firstDistance float32) {
	t.Helper()
	planes, err := file.Planes()
	if err != nil {
		t.Fatal(err)
	}
	if len(planes) != 2 || planes[0].Normal[0] != 1 || planes[0].Distance != firstDistance || planes[0].AxisType != 0 || planes[1].Normal[1] != 1 || planes[1].Distance != -9 || planes[1].AxisType != 1 {
		t.Fatalf("planes=%+v", planes)
	}
}

func TestCompressedOrdinaryReadEditAndPreservation(t *testing.T) {
	options := bsp.ReadOptions{Profile: bsp.ProfileSDK2013}
	file := readCompressedPlanes(t, options)
	requireTwoPlanes(t, file, 12)
	original, err := file.RawLump(1)
	if err != nil {
		t.Fatal(err)
	}
	logical, err := file.DecodedLump(1)
	if err != nil {
		t.Fatal(err)
	}
	logical[0] = 9
	requireTwoPlanes(t, file, 12)
	exported, err := file.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	reread, err := bsp.Read(bytes.NewReader(exported), options)
	if err != nil {
		t.Fatal(err)
	}
	stored, err := reread.RawLump(1)
	if err != nil || !bytes.Equal(stored, original) {
		t.Fatalf("stored=%x error=%v", stored, err)
	}
	requireTwoPlanes(t, reread, 12)
	planes, err := file.Planes()
	if err != nil {
		t.Fatal(err)
	}
	planes[0].AxisType = -1
	if err := file.SetPlanes(planes); !errors.Is(err, bsp.ErrMalformed) {
		t.Fatal(err)
	}
	unchanged, err := file.RawLump(1)
	if err != nil || !bytes.Equal(unchanged, original) {
		t.Fatalf("failed edit changed stored bytes: %x %v", unchanged, err)
	}
	planes[0].AxisType = 0
	planes[0].Distance = 31
	if err := file.SetPlanes(planes); err != nil {
		t.Fatal(err)
	}
	desc, err := file.Descriptor(1)
	if err != nil || desc.FourCC != [4]byte{} || desc.Length != 40 {
		t.Fatalf("metadata=%+v error=%v", desc, err)
	}
	edited, err := file.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	reread, err = bsp.Read(bytes.NewReader(edited), options)
	if err != nil {
		t.Fatal(err)
	}
	requireTwoPlanes(t, reread, 31)
}

func TestCompressionSelectionAndLimits(t *testing.T) {
	fixture := compressedPlaneFixture(t)
	unmarked, err := bsp.Read(bytes.NewReader(containerBytes(0, fixture)), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	if err != nil {
		t.Fatal(err)
	}
	text, err := unmarked.Entities()
	if err != nil || text != string(fixture) {
		t.Fatalf("unmarked prefix decoded: %x error=%v", text, err)
	}
	for _, options := range []bsp.ReadOptions{{Profile: bsp.ProfileSDK2013, MaxDecodedLumpBytes: 39}, {Profile: bsp.ProfileSDK2013, MaxLZMADictionaryBytes: 4095}} {
		file := readCompressedPlanes(t, options)
		if _, err := file.Planes(); !errors.Is(err, bsp.ErrLimitExceeded) {
			t.Fatal(err)
		}
		output, err := file.RawLump(1)
		if err != nil || !bytes.Equal(output, fixture) {
			t.Fatalf("limit changed raw=%x error=%v", output, err)
		}
	}
	requireTwoPlanes(t, readCompressedPlanes(t, bsp.ReadOptions{Profile: bsp.ProfileSDK2013, MaxDecodedLumpBytes: 40, MaxLZMADictionaryBytes: 4096}), 12)
	for _, options := range []bsp.ReadOptions{{Profile: bsp.ProfileSDK2013, MaxDecodedLumpBytes: -1}, {Profile: bsp.ProfileSDK2013, MaxLZMADictionaryBytes: -1}} {
		input := bytes.NewReader(containerBytes(0, []byte("unread")))
		if _, err := bsp.Read(input, options); !errors.Is(err, bsp.ErrInvalidOptions) {
			t.Fatal(err)
		}
		if input.Len() != bsp.HeaderSize+len("unread") {
			t.Fatal("invalid options consumed input")
		}
	}
}

func TestCompressedGameLogicalEditPreservesSiblingsAndFlags(t *testing.T) {
	file, err := bsp.Read(bytes.NewReader(containerBytes(-1, nil)), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	if err != nil {
		t.Fatal(err)
	}
	frame := compressedPlaneFixture(t)
	entries := []lumps.GameLump{{ID: 11, Flags: 0x8001, Version: 99, LogicalLength: 40, StoredData: frame}, {ID: 12, Flags: 0x20, Version: 77, StoredData: []byte("plain")}, {ID: 13, Flags: 1, Version: 88, LogicalLength: 40, StoredData: frame}, {}}
	if err := file.SetGameLumps(entries); err != nil {
		t.Fatal(err)
	}
	logical, err := file.DecodedGameLump(0)
	if err != nil {
		t.Fatal(err)
	}
	if len(logical) != 40 || binary.LittleEndian.Uint32(logical[12:]) != 0x41400000 {
		t.Fatalf("logical=%x", logical)
	}
	logical[0] = 9
	again, err := file.DecodedGameLump(0)
	if err != nil || again[0] != 0 {
		t.Fatalf("logical snapshot alias: %x %v", again, err)
	}
	if err := file.SetDecodedGameLump(0, []byte("edited")); err != nil {
		t.Fatal(err)
	}
	current, err := file.GameLumps()
	if err != nil {
		t.Fatal(err)
	}
	if len(current) != 4 || current[0].Flags != 0x8000 || current[0].ID != 11 || current[0].Version != 99 || current[0].LogicalLength != 6 || string(current[0].StoredData) != "edited" || string(current[1].StoredData) != "plain" || !bytes.Equal(current[2].StoredData, frame) {
		t.Fatalf("entries=%+v", current)
	}
	desc, err := file.Descriptor(35)
	if err != nil || desc.FourCC != [4]byte{} {
		t.Fatalf("outer=%+v error=%v", desc, err)
	}
	exported, err := file.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	reimport, err := bsp.Read(bytes.NewReader(exported), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	if err != nil {
		t.Fatal(err)
	}
	plain, err := reimport.DecodedGameLump(0)
	if err != nil || string(plain) != "edited" {
		t.Fatalf("edited=%x error=%v", plain, err)
	}
	sibling, err := reimport.DecodedGameLump(2)
	if err != nil || !bytes.Equal(sibling, again) {
		t.Fatalf("sibling=%x error=%v", sibling, err)
	}
}

func TestDecodedGameEditPreservesNonzeroGapsAndAlignment(t *testing.T) {
	frame := compressedPlaneFixture(t)
	raw := make([]byte, 116)
	binary.LittleEndian.PutUint32(raw, 3)
	putGameEntry(raw, 0, 11, 0x8001, 99, 1036+56, 40)
	putGameEntry(raw, 1, 12, 0x20, 77, 1036+108, 3)
	putGameEntry(raw, 2, 0, 0, 0, 1036+116, 0)
	copy(raw[52:], "HEAD")
	copy(raw[56:], frame)
	copy(raw[100:], "TAILGAP!")
	copy(raw[108:], "xyzAFTER")
	file := readGameFixture(t, raw)
	if err := file.SetDecodedGameLump(0, []byte("HELLO")); err != nil {
		t.Fatal(err)
	}
	for attempt := 0; attempt < 2; attempt++ {
		entries := requireGameEntries(t, file, 3)
		if entries[0].Flags != 0x8000 || entries[0].LogicalLength != 5 || string(entries[0].StoredData) != "HELLO" || string(entries[1].StoredData) != "xyz" {
			t.Fatalf("entries=%+v", entries)
		}
		for _, entry := range entries {
			if entry.Offset%4 != 0 {
				t.Fatalf("unaligned offset %d", entry.Offset)
			}
		}
		descriptor, _ := file.Descriptor(35)
		stored, err := file.RawLump(35)
		if err != nil {
			t.Fatal(err)
		}
		first := int(entries[0].Offset - descriptor.Offset)
		second := int(entries[1].Offset - descriptor.Offset)
		if string(stored[52:56]) != "HEAD" || string(stored[first+5:first+13]) != "TAILGAP!" || string(stored[second+3:second+8]) != "AFTER" {
			t.Fatalf("lost leading, compressed trailing or unclaimed sibling gap: %x", stored)
		}
		exported, err := file.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}
		file, err = bsp.Read(bytes.NewReader(exported), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
		if err != nil {
			t.Fatal(err)
		}
	}
}

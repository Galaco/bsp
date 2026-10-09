package bsp_test

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"os"
	"testing"

	"github.com/galaco/bsp"
)

// Expected map-state values were measured with unmodified Valve checksum_crc.cpp
// at b8cfb12c0e083a2ef5b2f9f9b50f3902fa034474, SDK numeric-ID traversal and
// CRC32_Init/ProcessBuffer without Final. See task14 for source hash and harness.
func controlledCRCWire(reverse bool) []byte {
	raw := make([]byte, bsp.HeaderSize+8)
	copy(raw, []byte("VBSP"))
	binary.LittleEndian.PutUint32(raw[4:], 20)
	first, last := bsp.HeaderSize, bsp.HeaderSize+4
	if reverse {
		first, last = last, first
	}
	binary.LittleEndian.PutUint32(raw[8+29*16:], uint32(first))
	binary.LittleEndian.PutUint32(raw[12+29*16:], 4)
	binary.LittleEndian.PutUint32(raw[16+29*16:], 77)
	binary.LittleEndian.PutUint32(raw[8+63*16:], uint32(last))
	binary.LittleEndian.PutUint32(raw[12+63*16:], 4)
	binary.LittleEndian.PutUint32(raw[16+63*16:], 99)
	copy(raw[first:first+4], "AAAA")
	copy(raw[last:last+4], "BBBB")
	return raw
}
func requireMapCRC(t *testing.T, file *bsp.File, want uint32) {
	t.Helper()
	got, err := file.MapCRC32()
	if err != nil || got != want {
		t.Fatalf("map CRC=%d error=%v want=%d", got, err, want)
	}
}
func TestMapCRCNumericOrderOpaqueLateIDAndExcludedBytes(t *testing.T) {
	for _, reverse := range []bool{false, true} {
		wire := controlledCRCWire(reverse)
		file, err := bsp.Read(bytes.NewReader(wire), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
		if err != nil {
			t.Fatal(err)
		}
		requireMapCRC(t, file, 614461499)
		if err := file.SetEntities("entities are excluded despite length changes"); err != nil {
			t.Fatal(err)
		}
		requireMapCRC(t, file, 614461499)
		binary.LittleEndian.PutUint32(wire[bsp.HeaderSize-4:], 999)
		wire = append(wire, []byte("ignored trailing bytes")...)
		other, err := bsp.Read(bytes.NewReader(wire), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
		if err != nil {
			t.Fatal(err)
		}
		requireMapCRC(t, other, 614461499)
		requireMapCRC(t, rereadOrdinary(t, file), 614461499)
		desc, err := file.Descriptor(bsp.LumpPhysCollide)
		if err != nil {
			t.Fatal(err)
		}
		if err := file.SetRawLump(bsp.LumpPhysCollide, desc, []byte("CCCC")); err != nil {
			t.Fatal(err)
		}
		requireMapCRC(t, file, 3324843012)
		// Reset first payload; late ID63 remains included even in an opaque profile.
		if err := file.SetRawLump(bsp.LumpPhysCollide, desc, []byte("AAAA")); err != nil {
			t.Fatal(err)
		}
		desc, err = file.Descriptor(bsp.LumpDispMultiBlend)
		if err != nil {
			t.Fatal(err)
		}
		if err := file.SetRawLump(bsp.LumpDispMultiBlend, desc, []byte("ZZZZ")); err != nil {
			t.Fatal(err)
		}
		requireMapCRC(t, file, 846232972)
	}
	// Independent physical gaps are not payload, even when nonzero.
	padded := make([]byte, 1056)
	copy(padded, controlledCRCWire(false)[:1040])
	for i := 1040; i < 1052; i++ {
		padded[i] = 0xa5
	}
	binary.LittleEndian.PutUint32(padded[8+63*16:], 1052)
	copy(padded[1052:], "BBBB")
	file, err := bsp.Read(bytes.NewReader(padded), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	if err != nil {
		t.Fatal(err)
	}
	requireMapCRC(t, file, 614461499)
	empty := spatialFixture(t, bsp.ProfileSDK2013, nil)
	requireMapCRC(t, empty, 0xffffffff)
}

func TestMapCRCStoredCompressionAndAtomicTypedTransitions(t *testing.T) {
	file := readCompressedPlanes(t, bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	requireMapCRC(t, file, 863388745)
	planes, err := file.Planes()
	if err != nil {
		t.Fatal(err)
	}
	requireMapCRC(t, file, 863388745)
	before, err := file.Descriptor(bsp.LumpPlanes)
	if err != nil {
		t.Fatal(err)
	}
	planes[0].AxisType = 99
	if err := file.SetPlanes(planes); !errors.Is(err, bsp.ErrMalformed) {
		t.Fatalf("invalid typed edit=%v", err)
	}
	requireMapCRC(t, file, 863388745)
	after, err := file.Descriptor(bsp.LumpPlanes)
	if err != nil || after != before {
		t.Fatal("failed edit changed compression metadata")
	}
	planes[0].AxisType = 0
	if err := file.SetPlanes(planes); err != nil {
		t.Fatal(err)
	}
	requireMapCRC(t, file, 3810161726)
	requireMapCRC(t, rereadOrdinary(t, file), 3810161726)
	planes[0].Distance = 13
	if err := file.SetPlanes(planes); err != nil {
		t.Fatal(err)
	}
	requireMapCRC(t, file, 1229467462)
	requireMapCRC(t, rereadOrdinary(t, file), 1229467462)
}

func TestMapCRCRelocatedGamePayloadAndPurity(t *testing.T) {
	wire := make([]byte, 2024)
	copy(wire, controlledCRCWire(false))
	binary.LittleEndian.PutUint32(wire[8+35*16:], 2000)
	binary.LittleEndian.PutUint32(wire[12+35*16:], 24)
	game := wire[2000:]
	binary.LittleEndian.PutUint32(game, 1)
	binary.LittleEndian.PutUint32(game[4:], 0x12345678)
	binary.LittleEndian.PutUint16(game[10:], 7)
	binary.LittleEndian.PutUint32(game[12:], 2020)
	binary.LittleEndian.PutUint32(game[16:], 4)
	copy(game[20:], "GAME")
	file, err := bsp.Read(bytes.NewReader(wire), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	if err != nil {
		t.Fatal(err)
	}
	header := file.Header()
	original, err := file.RawLump(bsp.LumpGame)
	if err != nil {
		t.Fatal(err)
	}
	requireMapCRC(t, file, 1289086378)
	requireMapCRC(t, file, 1289086378)
	if file.Header() != header {
		t.Fatal("CRC mutated source descriptors")
	}
	after, err := file.RawLump(bsp.LumpGame)
	if err != nil || !bytes.Equal(after, original) {
		t.Fatal("CRC mutated source game bytes")
	}
	next := rereadOrdinary(t, file)
	requireMapCRC(t, next, 1289086378)
	desc, err := next.Descriptor(bsp.LumpGame)
	if err != nil || desc.Offset != 1040 {
		t.Fatalf("planned game base=%+v %v", desc, err)
	}
	if err := file.SetRawLump(bsp.LumpGame, bsp.HeaderLump{Offset: 2000}, []byte{255, 255, 255, 255}); err != nil {
		t.Fatal(err)
	}
	if _, err := file.MapCRC32(); !errors.Is(err, bsp.ErrMalformed) {
		t.Fatalf("CRC hid invalid export state=%v", err)
	}
	if _, err := file.MarshalBinary(); !errors.Is(err, bsp.ErrMalformed) {
		t.Fatalf("writer accepted invalid export state=%v", err)
	}
}

func TestBundledMapCRCMatchesIndependentPlannedSDKReference(t *testing.T) {
	input, err := os.Open("ar_baggage.bsp.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()
	stream, err := gzip.NewReader(input)
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()
	file, err := bsp.Read(stream, bsp.ReadOptions{Profile: bsp.ProfileCSGO})
	if err != nil {
		t.Fatal(err)
	}
	// Old2836609078 was physical-order imported state. Independent numeric-ID
	// original state is4211417332; canonical game relocation produces2037425296.
	requireMapCRC(t, file, 2037425296)
	requireMapCRC(t, rereadOrdinary(t, file), 2037425296)
}

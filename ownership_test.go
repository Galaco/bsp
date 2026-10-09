package bsp_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"testing"

	"github.com/galaco/bsp"
)

func TestOwnedRawReplacement(t *testing.T) {
	input := containerBytes(0, []byte("old"))
	file, err := bsp.Read(bytes.NewReader(input), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	if err != nil {
		t.Fatal(err)
	}
	input[bsp.HeaderSize] = 'X'
	if text, err := file.Entities(); err != nil || text != "old" {
		t.Fatalf("text=%q err=%v", text, err)
	}
	desc, err := file.Descriptor(0)
	if err != nil {
		t.Fatal(err)
	}
	replacement := []byte("new entity")
	if err := file.SetRawLump(0, desc, replacement); err != nil {
		t.Fatal(err)
	}
	replacement[0] = 'X'
	raw, err := file.RawLump(0)
	if err != nil {
		t.Fatal(err)
	}
	raw[0] = 'Y'
	if text, err := file.Entities(); err != nil || text != "new entity" {
		t.Fatalf("text=%q err=%v", text, err)
	}
	desc, err = file.Descriptor(0)
	if err != nil || desc.Length != 10 {
		t.Fatalf("descriptor=%+v error=%v", desc, err)
	}
	header := file.Header()
	header.Lumps[0].Length = 99
	if file.Header().Lumps[0].Length != 10 {
		t.Fatal("header snapshot mutated file")
	}
}

func TestAtomicTypedReplacement(t *testing.T) {
	raw := make([]byte, 20)
	binary.LittleEndian.PutUint32(raw, math.Float32bits(1))
	binary.LittleEndian.PutUint32(raw[12:], math.Float32bits(7))
	file, err := bsp.Read(bytes.NewReader(containerBytes(1, raw)), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	if err != nil {
		t.Fatal(err)
	}
	planes, err := file.Planes()
	if err != nil {
		t.Fatal(err)
	}
	planes[0].Distance = 23
	unchanged, err := file.Planes()
	if err != nil || unchanged[0].Distance != 7 {
		t.Fatalf("snapshot mutation=%v error=%v", unchanged, err)
	}
	if err := file.SetPlanes(planes); err != nil {
		t.Fatal(err)
	}
	planes[0].Distance = 99
	current, err := file.Planes()
	if err != nil || current[0].Distance != 23 {
		t.Fatalf("committed=%v error=%v", current, err)
	}
	before, err := file.RawLump(1)
	if err != nil {
		t.Fatal(err)
	}
	current[0].AxisType = -1
	if err := file.SetPlanes(current); !errors.Is(err, bsp.ErrMalformed) {
		t.Fatal(err)
	}
	after, err := file.RawLump(1)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("failed edit changed bytes: %x error=%v", after, err)
	}
	if got := math.Float32frombits(binary.LittleEndian.Uint32(after[12:])); got != 23 {
		t.Fatalf("serialized distance=%g", got)
	}
}

func TestOpaqueSnapshotIsolation(t *testing.T) {
	file, err := bsp.Read(bytes.NewReader(containerBytes(63, []byte{1, 2, 3})), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	if err != nil {
		t.Fatal(err)
	}
	first, err := file.RawLump(63)
	if err != nil {
		t.Fatal(err)
	}
	first[0] = 9
	if err := file.SetRawLump(63, bsp.HeaderLump{Version: 987, FourCC: [4]byte{1, 2, 3, 4}}, []byte{4, 5}); err != nil {
		t.Fatal(err)
	}
	desc, err := file.Descriptor(63)
	if err != nil || desc.Version != 987 || desc.FourCC != [4]byte{1, 2, 3, 4} || desc.Length != 2 {
		t.Fatalf("metadata=%+v error=%v", desc, err)
	}
	stored, err := file.RawLump(63)
	if err != nil || !bytes.Equal(stored, []byte{4, 5}) {
		t.Fatalf("stored=%x error=%v", stored, err)
	}
}

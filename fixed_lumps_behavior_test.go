package bsp_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"testing"

	"github.com/galaco/bsp"
	"github.com/galaco/bsp/lumps"
)

func fixedAccessorFile(t *testing.T) *bsp.File {
	t.Helper()
	raw := make([]byte, bsp.HeaderSize)
	copy(raw, "VBSP")
	binary.LittleEndian.PutUint32(raw[4:], 20)
	file, err := bsp.Read(bytes.NewReader(raw), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	if err != nil {
		t.Fatal(err)
	}
	return file
}

func fixedAccessorReadback(t *testing.T, file *bsp.File) *bsp.File {
	t.Helper()
	raw, err := file.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	next, err := bsp.Read(bytes.NewReader(raw), bsp.ReadOptions{Profile: file.Profile()})
	if err != nil {
		t.Fatal(err)
	}
	return next
}

func TestFixedAccessorPersistentEditsPreservePaddingAndFloatBits(t *testing.T) {
	file := fixedAccessorFile(t)
	nodeWire := make([]byte, 32)
	binary.LittleEndian.PutUint16(nodeWire[30:], 0xa17e)
	if err := file.SetRawLump(bsp.LumpNodes, bsp.HeaderLump{}, nodeWire); err != nil {
		t.Fatal(err)
	}
	vertexWire := make([]byte, 12)
	bits := [3]uint32{0x7fc12345, 0x80000000, 0x3f800000}
	for i, value := range bits {
		binary.LittleEndian.PutUint32(vertexWire[i*4:], value)
	}
	if err := file.SetRawLump(bsp.LumpVertexes, bsp.HeaderLump{}, vertexWire); err != nil {
		t.Fatal(err)
	}
	nodes, err := file.Nodes()
	if err != nil {
		t.Fatal(err)
	}
	nodes[0].PlaneNum = 17
	untouched, err := file.Nodes()
	if err != nil || untouched[0].PlaneNum != 0 {
		t.Fatalf("getter edit changed stored nodes: %v, %v", untouched, err)
	}
	if err := file.SetNodes(nodes); err != nil {
		t.Fatal(err)
	}
	nodes[0].PlaneNum = 99
	vertices, err := file.Vertices()
	if err != nil {
		t.Fatal(err)
	}
	vertices[0][2] = 2
	if err := file.SetVertices(vertices); err != nil {
		t.Fatal(err)
	}
	vertices[0][0] = 0
	next := fixedAccessorReadback(t, file)
	gotNodes, err := next.Nodes()
	if err != nil || len(gotNodes) != 1 || gotNodes[0].PlaneNum != 17 || uint16(gotNodes[0].Padding) != 0xa17e {
		t.Fatalf("node edit/padding did not persist: %v, %v", gotNodes, err)
	}
	gotVertices, err := next.Vertices()
	if err != nil || len(gotVertices) != 1 {
		t.Fatalf("vertex readback: %v, %v", gotVertices, err)
	}
	wantBits := [3]uint32{bits[0], bits[1], math.Float32bits(2)}
	for i, want := range wantBits {
		if got := math.Float32bits(gotVertices[0][i]); got != want {
			t.Fatalf("vertex component %d bits = %#x, want %#x", i, got, want)
		}
	}
}

func TestFixedAccessorTextureBytesAndLightingChannelsPersist(t *testing.T) {
	file := fixedAccessorFile(t)
	text := "materials/a\x00\x00materials/b\x00"
	if err := file.SetTextureStringData(text); err != nil {
		t.Fatal(err)
	}
	ldr := []lumps.ColorRGBExponent32{{}}
	hdr := []lumps.ColorRGBExponent32{{}, {}}
	if err := file.SetLighting(ldr); err != nil {
		t.Fatal(err)
	}
	if err := file.SetLightingHDR(hdr); err != nil {
		t.Fatal(err)
	}
	next := fixedAccessorReadback(t, file)
	got, err := next.TextureStringData()
	if err != nil || got != text {
		t.Fatalf("texture bytes = %q, %v; want %q", got, err, text)
	}
	gotLDR, err := next.Lighting()
	if err != nil || len(gotLDR) != 1 {
		t.Fatalf("LDR channel = %v, %v", gotLDR, err)
	}
	gotHDR, err := next.LightingHDR()
	if err != nil || len(gotHDR) != 2 {
		t.Fatalf("HDR channel = %v, %v", gotHDR, err)
	}
}

func TestFixedAccessorUnknownRevisionFailureIsAtomic(t *testing.T) {
	file := fixedAccessorFile(t)
	original := []byte{1, 2, 3, 4}
	descriptor := bsp.HeaderLump{Version: 99, FourCC: [4]byte{9, 8, 7, 6}}
	if err := file.SetRawLump(bsp.LumpTexDataStringData, descriptor, original); err != nil {
		t.Fatal(err)
	}
	before, err := file.Descriptor(bsp.LumpTexDataStringData)
	if err != nil {
		t.Fatal(err)
	}
	if err := file.SetTextureStringData("replacement\x00"); !errors.Is(err, bsp.ErrUnsupportedLump) {
		t.Fatalf("unknown revision setter = %v", err)
	}
	after, err := file.Descriptor(bsp.LumpTexDataStringData)
	if err != nil || before != after {
		t.Fatalf("failed setter changed descriptor: %v -> %v, %v", before, after, err)
	}
	raw, err := file.RawLump(bsp.LumpTexDataStringData)
	if err != nil || !bytes.Equal(raw, original) {
		t.Fatalf("failed setter changed raw: %v, %v", raw, err)
	}
	if err := file.SetRawLump(bsp.LumpTexDataStringData, bsp.HeaderLump{FourCC: descriptor.FourCC}, original); err != nil {
		t.Fatal(err)
	}
	if err := file.SetTextureStringData("replacement\x00"); err != nil {
		t.Fatal(err)
	}
	after, err = file.Descriptor(bsp.LumpTexDataStringData)
	if err != nil || after.Version != 0 || after.FourCC != [4]byte{} || after.Length != int32(len("replacement\x00")) {
		t.Fatalf("successful setter metadata = %v, %v", after, err)
	}
}

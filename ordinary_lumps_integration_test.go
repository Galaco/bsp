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

type ordinaryWire struct {
	version int32
	data    []byte
}

func readOrdinaryFixture(t *testing.T, profile bsp.Profile, payloads map[int]ordinaryWire) *bsp.File {
	t.Helper()
	raw := make([]byte, bsp.HeaderSize)
	copy(raw, "VBSP")
	version := uint32(20)
	if profile != bsp.ProfileSDK2013 {
		version = 21
	}
	binary.LittleEndian.PutUint32(raw[4:], version)
	for id := 0; id < bsp.LumpCount; id++ {
		wire, ok := payloads[id]
		if !ok {
			continue
		}
		offset := len(raw)
		raw = append(raw, wire.data...)
		entry := raw[8+id*16 : 8+(id+1)*16]
		if profile == bsp.ProfileL4D2 {
			binary.LittleEndian.PutUint32(entry, uint32(wire.version))
			binary.LittleEndian.PutUint32(entry[4:], uint32(offset))
			binary.LittleEndian.PutUint32(entry[8:], uint32(len(wire.data)))
		} else {
			binary.LittleEndian.PutUint32(entry, uint32(offset))
			binary.LittleEndian.PutUint32(entry[4:], uint32(len(wire.data)))
			binary.LittleEndian.PutUint32(entry[8:], uint32(wire.version))
		}
	}
	file, err := bsp.Read(bytes.NewReader(raw), bsp.ReadOptions{Profile: profile})
	if err != nil {
		t.Fatal(err)
	}
	return file
}
func rereadOrdinary(t *testing.T, file *bsp.File) *bsp.File {
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
func twoLeafWire(version int32) []byte {
	width := 32
	if version == 0 {
		width = 56
	}
	raw := make([]byte, width*2)
	for i := 0; i < 2; i++ {
		p := raw[i*width:]
		binary.LittleEndian.PutUint32(p, uint32(11+i))
		binary.LittleEndian.PutUint16(p[4:], uint16(i))
		binary.LittleEndian.PutUint16(p[6:], uint16(0xfe00+i))
		binary.LittleEndian.PutUint16(p[28:], 0xffff)
		p[width-2] = byte(0xa1 + i)
		p[width-1] = byte(0xc1 + i)
		if version == 0 {
			p[30] = byte(71 + i)
			p[53] = 0x80
		}
	}
	return raw
}
func twoLightWire() []byte {
	raw := make([]byte, 200)
	for i := 0; i < 2; i++ {
		p := raw[i*100:]
		for axis := 0; axis < 3; axis++ {
			binary.LittleEndian.PutUint32(p[axis*4:], math.Float32bits(float32(10+i+axis)))
			binary.LittleEndian.PutUint32(p[36+axis*4:], math.Float32bits(float32(1+i+axis)))
		}
		binary.LittleEndian.PutUint32(p[48:], uint32(2+i))
		binary.LittleEndian.PutUint32(p[52:], uint32(1+i))
		binary.LittleEndian.PutUint32(p[56:], uint32(17+i))
		binary.LittleEndian.PutUint32(p[96:], uint32(40+i))
	}
	return raw
}

func TestOrdinaryPublicLayoutsSnapshotsAndEdits(t *testing.T) {
	for _, leafVersion := range []int32{0, 1} {
		file := readOrdinaryFixture(t, bsp.ProfileCSGO, map[int]ordinaryWire{10: {leafVersion, twoLeafWire(leafVersion)}, 15: {1, twoLightWire()}, 54: {1, twoLightWire()}, 19: {0, []byte{1, 0, 2, 0, 3, 0, 1, 2, 4, 0, 5, 0, 6, 0, 7, 8}}})
		leaves, err := file.Leaves()
		if err != nil {
			t.Fatal(err)
		}
		if len(leaves) != 2 || leaves[1].Contents != 12 || leaves[1].Padding != [2]byte{0xa2, 0xc2} || leaves[0].HasInlineLighting != (leafVersion == 0) {
			t.Fatalf("leaves=%+v", leaves)
		}
		leaves[0].Contents = 51
		mustLeaves, err := file.Leaves()
		if err != nil || mustLeaves[0].Contents != 11 {
			t.Fatalf("leaf snapshot alias: %+v %v", mustLeaves, err)
		}
		if err := file.SetLeaves(leaves); err != nil {
			t.Fatal(err)
		}
		ldr, err := file.WorldLights()
		if err != nil {
			t.Fatal(err)
		}
		hdr, err := file.WorldLightsHDR()
		if err != nil {
			t.Fatal(err)
		}
		if len(ldr) != 2 || ldr[1].ShadowCastOffset[2] != 4 || ldr[1].Type != lumps.EmitSpotLight || hdr[1].Owner != 41 {
			t.Fatalf("lights=%+v HDR=%+v", ldr, hdr)
		}
		ldr[0].ShadowCastOffset[1] = 9
		hdr[1].Owner = 77
		if err := file.SetWorldLights(ldr); err != nil {
			t.Fatal(err)
		}
		if err := file.SetWorldLightsHDR(hdr); err != nil {
			t.Fatal(err)
		}
		sides, err := file.BrushSides()
		if err != nil {
			t.Fatal(err)
		}
		if len(sides) != 2 || !sides[0].HasThin || sides[0].Bevel != 1 || sides[0].Thin != 2 || sides[1].Thin != 8 {
			t.Fatalf("sides=%+v", sides)
		}
		sides[1].Thin = 21
		if err := file.SetBrushSides(sides); err != nil {
			t.Fatal(err)
		}
		file = rereadOrdinary(t, file)
		leaves, err = file.Leaves()
		if err != nil || leaves[0].Contents != 51 || leaves[1].Padding != [2]byte{0xa2, 0xc2} {
			t.Fatalf("edited leaves=%+v %v", leaves, err)
		}
		ldr, err = file.WorldLights()
		if err != nil || ldr[0].ShadowCastOffset[1] != 9 {
			t.Fatalf("edited LDR=%+v %v", ldr, err)
		}
		hdr, err = file.WorldLightsHDR()
		if err != nil || hdr[1].Owner != 77 {
			t.Fatalf("edited HDR=%+v %v", hdr, err)
		}
		sides, err = file.BrushSides()
		if err != nil || sides[1].Thin != 21 {
			t.Fatalf("edited sides=%+v %v", sides, err)
		}
		for _, id := range []bsp.LumpID{15, 54} {
			desc, err := file.Descriptor(id)
			if err != nil || desc.Version != 1 || desc.Length != 200 {
				t.Fatalf("descriptor=%+v %v", desc, err)
			}
		}
	}
}
func TestSameSizeBrushSideDialectsAndUnsupportedWorldlight(t *testing.T) {
	raw := []byte{1, 0, 2, 0, 3, 0, 1, 2}
	for _, profile := range []bsp.Profile{bsp.ProfileSDK2013, bsp.ProfileL4D2, bsp.ProfileCSGO} {
		file := readOrdinaryFixture(t, profile, map[int]ordinaryWire{19: {0, raw}})
		sides, err := file.BrushSides()
		if err != nil {
			t.Fatal(err)
		}
		if profile == bsp.ProfileCSGO {
			if sides[0].Bevel != 1 || sides[0].Thin != 2 || !sides[0].HasThin {
				t.Fatal(sides)
			}
		} else {
			if sides[0].Bevel != 513 || sides[0].Thin != 0 || sides[0].HasThin {
				t.Fatal(sides)
			}
			before, _ := file.RawLump(19)
			sides[0].Thin = 1
			if err := file.SetBrushSides(sides); err == nil {
				t.Fatal("accepted unrepresentable thin field")
			}
			after, _ := file.RawLump(19)
			if !bytes.Equal(before, after) {
				t.Fatal("failed brush setter committed")
			}
		}
	}
	for _, profile := range []bsp.Profile{bsp.ProfileSDK2013, bsp.ProfileL4D2} {
		file := readOrdinaryFixture(t, profile, map[int]ordinaryWire{15: {1, twoLightWire()}})
		if _, err := file.WorldLights(); !errors.Is(err, bsp.ErrUnsupportedLump) {
			t.Fatalf("unsupported %s: %v", profile, err)
		}
		file = rereadOrdinary(t, file)
		raw, err := file.RawLump(15)
		if err != nil || !bytes.Equal(raw, twoLightWire()) {
			t.Fatalf("opaque worldlight changed: %v", err)
		}
	}
}

func TestAmbientPublicCrossLumpValidationAndAtomicEdits(t *testing.T) {
	indices := []byte{2, 0, 0, 0, 0, 0, 0, 0}
	samples := make([]byte, 56)
	samples[0] = 11
	samples[24] = 1
	samples[27] = 0xab
	samples[28] = 22
	samples[52] = 9
	samples[55] = 0xcd
	file := readOrdinaryFixture(t, bsp.ProfileSDK2013, map[int]ordinaryWire{10: {1, twoLeafWire(1)}, 52: {0, indices}, 56: {1, samples}})
	data, err := file.AmbientSamples()
	if err != nil || len(data) != 2 || !data[0].HasPosition || data[1].Cube.Color[0].R != 22 || data[1].X != 9 || data[1].Pad != 0xcd {
		t.Fatalf("samples=%+v %v", data, err)
	}
	data[1].X = 31
	if err := file.SetAmbientSamples(data); err != nil {
		t.Fatal(err)
	}
	index, err := file.AmbientIndices()
	if err != nil || index[1].AmbientSampleCount != 0 || index[1].FirstAmbientSample != 0 {
		t.Fatalf("indices=%+v %v", index, err)
	}
	before, _ := file.RawLump(52)
	index[1].FirstAmbientSample = 2
	if err := file.SetAmbientIndices(index); err == nil {
		t.Fatal("out-of-range leaf-copy reference accepted")
	}
	after, _ := file.RawLump(52)
	if !bytes.Equal(before, after) {
		t.Fatal("failed index setter committed")
	}
	leaves, err := file.Leaves()
	if err != nil {
		t.Fatal(err)
	}
	leafBefore, _ := file.RawLump(10)
	if err := file.SetLeaves(leaves[:1]); err == nil {
		t.Fatal("accepted mismatched ambient leaf count")
	}
	leafAfter, _ := file.RawLump(10)
	if !bytes.Equal(leafBefore, leafAfter) {
		t.Fatal("failed leaf setter committed")
	}
	file = rereadOrdinary(t, file)
	data, err = file.AmbientSamples()
	if err != nil || data[1].X != 31 || data[1].Pad != 0xcd {
		t.Fatalf("edited samples=%+v %v", data, err)
	}
	hdr, err := file.AmbientSamplesHDR()
	if err != nil || len(hdr) != 0 {
		t.Fatalf("absent HDR gained fallback: %+v %v", hdr, err)
	}
}
func TestAmbientLegacyVersionsAndNoImplicitMigration(t *testing.T) {
	raw := make([]byte, 48)
	raw[0] = 41
	raw[24] = 42
	file := readOrdinaryFixture(t, bsp.ProfileSDK2013, map[int]ordinaryWire{10: {1, twoLeafWire(1)}, 55: {777, raw}, 56: {999, raw}})
	for _, get := range []func() ([]lumps.AmbientSample, error){file.AmbientSamples, file.AmbientSamplesHDR} {
		samples, err := get()
		if err != nil || len(samples) != 2 || samples[0].HasPosition || samples[1].Cube.Color[0].R != 42 {
			t.Fatalf("legacy=%+v %v", samples, err)
		}
	}
	before, _ := file.RawLump(56)
	if err := file.SetAmbientIndices([]lumps.AmbientIndex{{AmbientSampleCount: 1}, {AmbientSampleCount: 1, FirstAmbientSample: 1}}); err == nil {
		t.Fatal("implicitly migrated legacy samples")
	}
	after, _ := file.RawLump(56)
	desc, _ := file.Descriptor(52)
	if !bytes.Equal(before, after) || desc.Length != 0 {
		t.Fatal("failed migration committed")
	}
	samples, err := file.AmbientSamples()
	if err != nil {
		t.Fatal(err)
	}
	samples[1].Cube.Color[0].R = 91
	if err := file.SetAmbientSamples(samples); err != nil {
		t.Fatal(err)
	}
	file = rereadOrdinary(t, file)
	samples, err = file.AmbientSamples()
	if err != nil || samples[1].Cube.Color[0].R != 91 || samples[1].HasPosition {
		t.Fatalf("edited legacy=%+v %v", samples, err)
	}
	desc, _ = file.Descriptor(56)
	if desc.Version != 999 || desc.Length != 48 {
		t.Fatalf("legacy metadata=%+v", desc)
	}
}

func occlusionWire(version int32) []byte {
	width := 36
	if version == 2 {
		width = 40
	}
	raw := make([]byte, 12+width*2+12*2+4*6)
	binary.LittleEndian.PutUint32(raw, 2)
	at := 4
	for i := 0; i < 2; i++ {
		p := raw[at:]
		binary.LittleEndian.PutUint32(p, uint32(13+i))
		binary.LittleEndian.PutUint32(p[4:], uint32(i))
		binary.LittleEndian.PutUint32(p[8:], 1)
		binary.LittleEndian.PutUint32(p[12:], math.Float32bits(float32(21+i)))
		if version == 2 {
			binary.LittleEndian.PutUint32(p[36:], uint32(2+i))
		}
		at += width
	}
	binary.LittleEndian.PutUint32(raw[at:], 2)
	at += 4
	for i := 0; i < 2; i++ {
		binary.LittleEndian.PutUint32(raw[at:], uint32(i*3))
		binary.LittleEndian.PutUint32(raw[at+4:], 3)
		binary.LittleEndian.PutUint32(raw[at+8:], uint32(i))
		at += 12
	}
	binary.LittleEndian.PutUint32(raw[at:], 6)
	at += 4
	for i, index := range []uint32{0, 1, 2, 1, 2, 3} {
		binary.LittleEndian.PutUint32(raw[at+i*4:], index)
	}
	return raw
}
func TestOcclusionPublicReferencesAndAtomicReplacement(t *testing.T) {
	for _, version := range []int32{1, 2} {
		file := readOrdinaryFixture(t, bsp.ProfileSDK2013, map[int]ordinaryWire{1: {0, make([]byte, 40)}, 3: {0, make([]byte, 48)}, 9: {version, occlusionWire(version)}})
		value, err := file.Occlusion()
		if err != nil || len(value.Data) != 2 || value.Data[1].Flags != 14 || value.Polygons[1].PlaneNum != 1 || value.VertexIndices[5] != 3 || value.Data[0].HasArea != (version == 2) {
			t.Fatalf("occlusion=%+v %v", value, err)
		}
		before, _ := file.RawLump(9)
		value.Polygons[1].PlaneNum = 2
		if err := file.SetOcclusion(value); err == nil {
			t.Fatal("out-of-range occlusion plane accepted")
		}
		after, _ := file.RawLump(9)
		if !bytes.Equal(before, after) {
			t.Fatal("failed occlusion setter committed")
		}
		value.Polygons[1].PlaneNum = 1
		value.Data[1].Flags = 37
		if err := file.SetOcclusion(value); err != nil {
			t.Fatal(err)
		}
		file = rereadOrdinary(t, file)
		value, err = file.Occlusion()
		if err != nil || value.Data[1].Flags != 37 || value.VertexIndices[5] != 3 {
			t.Fatalf("edited occlusion=%+v %v", value, err)
		}
	}
}

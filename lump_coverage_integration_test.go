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

func displacementPayloads() map[bsp.LumpID][]byte {
	info := make([]byte, 2*176)
	faces := make([]byte, 2*56)
	for i := 0; i < 2; i++ {
		record := info[i*176 : (i+1)*176]
		binary.LittleEndian.PutUint32(record, uint32(math.Float32bits(float32(i*100+3))))
		binary.LittleEndian.PutUint32(record[12:], uint32(i*289))
		binary.LittleEndian.PutUint32(record[16:], uint32(i*512))
		binary.LittleEndian.PutUint32(record[20:], 4)
		binary.LittleEndian.PutUint32(record[24:], lumps.DisplacementFlagMagic|0x1234)
		binary.LittleEndian.PutUint16(record[36:], uint16(i))
		record[38], record[39] = byte(0xa1+i), 0xb2
		binary.LittleEndian.PutUint32(record[40:], 0xfffffffe)
		if i == 1 {
			binary.LittleEndian.PutUint32(record[44:], 17)
		}
		for edge := 0; edge < 8; edge++ {
			at := 48 + 6*edge
			binary.LittleEndian.PutUint16(record[at:], 0xffff)
			copy(record[at+2:at+6], []byte{0xe1, 0xe2, 0xe3, byte(edge + 1)})
		}
		binary.LittleEndian.PutUint16(record[48:], uint16(1-i))
		copy(record[50:54], []byte{3, 2, 1, 0x9a})
		for corner := 0; corner < 4; corner++ {
			at := 96 + 10*corner
			for n := 0; n < 4; n++ {
				binary.LittleEndian.PutUint16(record[at+n*2:], 0xf000+uint16(n))
			}
			record[at+9] = 0xc3
		}
		binary.LittleEndian.PutUint16(record[96:], 0)
		binary.LittleEndian.PutUint16(record[98:], 1)
		record[104] = 2
		for word := 0; word < 10; word++ {
			binary.LittleEndian.PutUint32(record[136+word*4:], uint32(0x10000+i*100+word))
		}
		binary.LittleEndian.PutUint32(faces[i*56+36:], 2) // three by one samples, regardless of styles.
		copy(faces[i*56+16:i*56+20], []byte{0, 1, 2, 3})
	}
	samples := []byte{254, 11, 12, 13, 255, 0, 21, 22, 23, 255, 255, 31, 32, 33, 0xde, 0xad, 0xbe, 0, 41, 42, 43, 1, 51, 52, 53, 2, 61, 62, 63, 0xfa, 0xfb}
	return map[bsp.LumpID][]byte{bsp.LumpDispInfo: info, bsp.LumpFaces: faces, bsp.LumpDispVerts: make([]byte, 578*20), bsp.LumpDispTris: make([]byte, 1024*2), bsp.LumpDispLightmapSamplePositions: samples}
}

func TestCompleteDisplacementFieldsAndGeometryAtomicity(t *testing.T) {
	for _, profile := range []bsp.Profile{bsp.ProfileSDK2013, bsp.ProfileCSGO, bsp.ProfileL4D2, bsp.ProfileReactiveDrop} {
		t.Run(string(profile), func(t *testing.T) {
			payloads := displacementPayloads()
			file := spatialFixture(t, profile, payloads)
			records, err := file.Displacements()
			if err != nil {
				t.Fatal(err)
			}
			if len(records) != 2 || records[1].MapFacePadding != [2]byte{0xa2, 0xb2} || records[1].AllowedVerts[9] != 0x1006d || records[1].Flags != lumps.DisplacementFlagMagic|0x1234 || records[1].MinTess != 0 || records[1].EdgeNeighbors[0].SubNeighbors[1].Orientation != 0xe1 || records[1].CornerNeighbors[0].Neighbors[3] != 0xf003 {
				t.Fatalf("complete fields=%+v", records)
			}
			before, _ := file.RawLump(bsp.LumpDispInfo)
			descriptor, _ := file.Descriptor(bsp.LumpDispInfo)
			records[1].MapFace = 2
			if err := file.SetDisplacements(records); !errors.Is(err, bsp.ErrMalformed) {
				t.Fatalf("invalid face error=%v", err)
			}
			raw, _ := file.RawLump(bsp.LumpDispInfo)
			after, _ := file.Descriptor(bsp.LumpDispInfo)
			if !bytes.Equal(raw, before) || after != descriptor {
				t.Fatal("invalid geometry partially committed")
			}
			records[1].MapFace = 1
			records[1].AllowedVerts[9] = 0xfedcba98
			records[1].CornerNeighbors[3].Padding = 0x57
			records[1].EdgeNeighbors[0].SubNeighbors[0].Padding = 0x68
			fresh, err := file.Displacements()
			if err != nil || fresh[1].AllowedVerts[9] != 0x1006d {
				t.Fatalf("detached mutation leaked: %v", err)
			}
			if err := file.SetDisplacements(records); err != nil {
				t.Fatal(err)
			}
			binary.LittleEndian.PutUint32(payloads[bsp.LumpDispInfo][176+172:], 0xfedcba98)
			payloads[bsp.LumpDispInfo][176+135] = 0x57
			payloads[bsp.LumpDispInfo][176+53] = 0x68
			next := rereadOrdinary(t, file)
			raw, _ = next.RawLump(bsp.LumpDispInfo)
			if !bytes.Equal(raw, payloads[bsp.LumpDispInfo]) {
				t.Fatal("displacement tail/padding bytes lost on export")
			}
			got, err := next.Displacements()
			if err != nil || got[1].AllowedVerts[9] != 0xfedcba98 || got[1].CornerNeighbors[3].Padding != 0x57 {
				t.Fatalf("displacement reread=%+v %v", got, err)
			}
		})
	}
}

func TestDisplacementSampleSpliceRebasesOnlyFollowingOffsets(t *testing.T) {
	payloads := displacementPayloads()
	file := spatialFixture(t, bsp.ProfileCSGO, payloads)
	values, err := file.DisplacementSamples(0)
	if err != nil {
		t.Fatal(err)
	}
	if len(values) != 3 || values[0].Index != 254 || values[1].Index != 255 || values[2].Index != 510 || values[2].Barycentric != [3]byte{31, 32, 33} {
		t.Fatalf("variable samples=%+v", values)
	}
	nextValues, err := file.DisplacementSamples(1)
	if err != nil || len(nextValues) != 3 || nextValues[2].Index != 2 {
		t.Fatalf("following samples=%+v %v", nextValues, err)
	}
	originalInfo, _ := file.RawLump(bsp.LumpDispInfo)
	originalSamples, _ := file.RawLump(bsp.LumpDispLightmapSamplePositions)
	header := file.Header()
	bad := append([]lumps.DisplacementSample(nil), values...)
	bad[1].Index = 511
	if err := file.SetDisplacementSamples(0, bad); !errors.Is(err, bsp.ErrMalformed) {
		t.Fatalf("unrepresentable index error=%v", err)
	}
	info, _ := file.RawLump(bsp.LumpDispInfo)
	samples, _ := file.RawLump(bsp.LumpDispLightmapSamplePositions)
	if !bytes.Equal(info, originalInfo) || !bytes.Equal(samples, originalSamples) || header != file.Header() {
		t.Fatal("failed stream encode changed state")
	}
	values[0].Index = 510
	values[0].Barycentric = [3]byte{201, 202, 203}
	if err := file.SetDisplacementSamples(0, values); err != nil {
		t.Fatal(err)
	}
	expected := append([]byte{255, 255, 201, 202, 203}, originalSamples[4:]...)
	binary.LittleEndian.PutUint32(originalInfo[176+44:], 18)
	next := rereadOrdinary(t, file)
	samples, _ = next.RawLump(bsp.LumpDispLightmapSamplePositions)
	info, _ = next.RawLump(bsp.LumpDispInfo)
	if !bytes.Equal(samples, expected) || !bytes.Equal(info, originalInfo) {
		t.Fatal("splice changed gaps, trailer, other fields or following offset")
	}
	rereadValues, err := next.DisplacementSamples(1)
	if err != nil || len(rereadValues) != 3 || rereadValues[0] != nextValues[0] || rereadValues[2] != nextValues[2] {
		t.Fatalf("rebased samples=%+v %v", rereadValues, err)
	}
	before, _ := next.RawLump(bsp.LumpDispInfo)
	beforeDesc, _ := next.Descriptor(bsp.LumpDispInfo)
	values, err = next.DisplacementSamples(0)
	if err != nil {
		t.Fatal(err)
	}
	values[1].Barycentric[2] = 99
	if err := next.SetDisplacementSamples(0, values); err != nil {
		t.Fatal(err)
	}
	info, _ = next.RawLump(bsp.LumpDispInfo)
	afterDesc, _ := next.Descriptor(bsp.LumpDispInfo)
	if !bytes.Equal(info, before) || beforeDesc != afterDesc {
		t.Fatal("same-width sample edit unnecessarily changed lump26")
	}
	if _, err := next.DisplacementSamples(-1); !errors.Is(err, bsp.ErrInvalidDisplacementIndex) {
		t.Fatalf("bad index=%v", err)
	}
}

func TestDisplacementSampleTruncationAndSharedSpanAreAtomic(t *testing.T) {
	for _, shared := range []bool{false, true} {
		t.Run(map[bool]string{false: "truncated", true: "shared"}[shared], func(t *testing.T) {
			payloads := displacementPayloads()
			if shared {
				binary.LittleEndian.PutUint32(payloads[bsp.LumpDispInfo][176+44:], 0)
			} else {
				payloads[bsp.LumpDispLightmapSamplePositions] = []byte{255, 0, 1, 2}
			}
			file := spatialFixture(t, bsp.ProfileCSGO, payloads)
			before := file.Header()
			raw, _ := file.RawLump(bsp.LumpDispLightmapSamplePositions)
			err := file.SetDisplacementSamples(0, []lumps.DisplacementSample{{Index: 254}, {Index: 255}, {Index: 510}})
			sentinel := bsp.ErrMalformed
			if shared {
				sentinel = bsp.ErrUnsupportedLump
			}
			if !errors.Is(err, sentinel) {
				t.Fatalf("invalid stream error=%v", err)
			}
			after, _ := file.RawLump(bsp.LumpDispLightmapSamplePositions)
			if before != file.Header() || !bytes.Equal(raw, after) {
				t.Fatal("failed splice changed authoritative state")
			}
		})
	}
}

func TestMultiblendVertexOwnershipAndLateColorEdit(t *testing.T) {
	payloads := displacementPayloads()
	binary.LittleEndian.PutUint32(payloads[bsp.LumpDispInfo][24:], lumps.DisplacementFlagMagic|lumps.DisplacementFlagHasMultiblend)
	blend := make([]byte, 289*80)
	for i := 0; i < 289; i++ {
		for j := 0; j < 20; j++ {
			binary.LittleEndian.PutUint32(blend[i*80+j*4:], math.Float32bits(float32(i*100+j)))
		}
	}
	payloads[bsp.LumpDispMultiBlend] = blend
	file := spatialFixture(t, bsp.ProfileCSGO, payloads)
	values, err := file.DisplacementMultiblends()
	if err != nil || len(values) != 289 || values[288].Colors[3][2] != 28819 {
		t.Fatalf("multiblend=%v len%d %v", values, len(values), err)
	}
	original, _ := file.RawLump(bsp.LumpDispMultiBlend)
	if err := file.SetDisplacementMultiblends(values[:288]); !errors.Is(err, bsp.ErrMalformed) {
		t.Fatalf("missing vertex record=%v", err)
	}
	raw, _ := file.RawLump(bsp.LumpDispMultiBlend)
	if !bytes.Equal(raw, original) {
		t.Fatal("bad multiblend count committed")
	}
	values[288].Colors[3][2] = -7.25
	if err := file.SetDisplacementMultiblends(values); err != nil {
		t.Fatal(err)
	}
	binary.LittleEndian.PutUint32(original[288*80+76:], math.Float32bits(-7.25))
	next := rereadOrdinary(t, file)
	raw, _ = next.RawLump(bsp.LumpDispMultiBlend)
	if !bytes.Equal(raw, original) {
		t.Fatal("late color or other components lost")
	}
	if err := next.SetRawLump(bsp.LumpDispMultiBlend, bsp.HeaderLump{}, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := next.DisplacementMultiblends(); !errors.Is(err, bsp.ErrMalformed) {
		t.Fatalf("missing flagged multiblend=%v", err)
	}
}

func TestWaterOverlayLastFaceAndSystemLevelOwnership(t *testing.T) {
	for _, profile := range []bsp.Profile{bsp.ProfileSDK2013, bsp.ProfileCSGO, bsp.ProfileL4D2, bsp.ProfileReactiveDrop} {
		t.Run(string(profile), func(t *testing.T) {
			water := make([]byte, 2240)
			for i := 0; i < 2; i++ {
				record := water[i*1120 : (i+1)*1120]
				binary.LittleEndian.PutUint32(record, uint32(40+i))
				binary.LittleEndian.PutUint16(record[6:], 3<<14|256)
				for j := 0; j < 256; j++ {
					binary.LittleEndian.PutUint32(record[8+j*4:], uint32(j))
				}
				for at := 1032; at < 1120; at += 4 {
					binary.LittleEndian.PutUint32(record[at:], math.Float32bits(float32(at+i*2000)))
				}
			}
			payloads := map[bsp.LumpID][]byte{bsp.LumpWaterOverlays: water, bsp.LumpFaces: make([]byte, 256*56), bsp.LumpTexInfo: make([]byte, 72)}
			if profile != bsp.ProfileSDK2013 {
				payloads[bsp.LumpOverlays] = make([]byte, 2*352)
				payloads[bsp.LumpOverlaySystemLevels] = []byte{1, 2, 3, 4, 21, 22, 23, 24}
			}
			file := spatialFixture(t, profile, payloads)
			values, err := file.WaterOverlays()
			if err != nil || len(values) != 2 || values[1].Faces[255] != 255 || values[1].BasisNormal[2] != 3116 || values[0].RenderOrder() != 3 {
				t.Fatalf("water values=%+v %v", values, err)
			}
			values[1].Faces[255] = 254
			values[1].BasisNormal[2] = -9.5
			if err := file.SetWaterOverlays(values); err != nil {
				t.Fatal(err)
			}
			binary.LittleEndian.PutUint32(water[1120+1028:], 254)
			binary.LittleEndian.PutUint32(water[1120+1116:], math.Float32bits(-9.5))
			if profile == bsp.ProfileSDK2013 {
				if _, err := file.OverlaySystemLevels(); !errors.Is(err, bsp.ErrUnsupportedLump) {
					t.Fatalf("SDK alias admitted: %v", err)
				}
			} else {
				levels, err := file.OverlaySystemLevels()
				if err != nil || len(levels) != 2 || levels[1].MaxGPULevel != 24 {
					t.Fatalf("system levels=%+v %v", levels, err)
				}
				before, _ := file.RawLump(bsp.LumpOverlays)
				if err := file.SetOverlays(nil); !errors.Is(err, bsp.ErrMalformed) {
					t.Fatalf("dangling system levels error=%v", err)
				}
				after, _ := file.RawLump(bsp.LumpOverlays)
				if !bytes.Equal(before, after) {
					t.Fatal("invalid companion count committed")
				}
				levels[1].MaxGPULevel = 231
				if err := file.SetOverlaySystemLevels(levels); err != nil {
					t.Fatal(err)
				}
				payloads[bsp.LumpOverlaySystemLevels][7] = 231
			}
			next := rereadOrdinary(t, file)
			raw, _ := next.RawLump(bsp.LumpWaterOverlays)
			if !bytes.Equal(raw, water) {
				t.Fatal("water face255 or geometry lost")
			}
			if profile != bsp.ProfileSDK2013 {
				raw, _ := next.RawLump(bsp.LumpOverlaySystemLevels)
				if !bytes.Equal(raw, payloads[bsp.LumpOverlaySystemLevels]) {
					t.Fatal("late system level edit lost")
				}
			}
		})
	}
}

func TestReactiveDropCollisionWireEditsWithoutInventedBlobUnits(t *testing.T) {
	payloads := map[bsp.LumpID][]byte{bsp.LumpPropCollision: make([]byte, 16), bsp.LumpPropHulls: make([]byte, 32), bsp.LumpPropHullVerts: make([]byte, 24), bsp.LumpPropTris: make([]byte, 16), bsp.LumpPropBlob: {0xa1, 0, 0xff}}
	for i := 0; i < 2; i++ {
		binary.LittleEndian.PutUint32(payloads[bsp.LumpPropCollision][i*8:], uint32(3+i))
		binary.LittleEndian.PutUint32(payloads[bsp.LumpPropCollision][i*8+4:], uint32(7+i))
		binary.LittleEndian.PutUint32(payloads[bsp.LumpPropHulls][i*16:], uint32(5+i))
		binary.LittleEndian.PutUint32(payloads[bsp.LumpPropHulls][i*16+4:], uint32(11+i))
		binary.LittleEndian.PutUint32(payloads[bsp.LumpPropHulls][i*16+8:], uint32(0xfffffffa+i))
		binary.LittleEndian.PutUint32(payloads[bsp.LumpPropHulls][i*16+12:], 0x80000001)
		for j := 0; j < 3; j++ {
			binary.LittleEndian.PutUint32(payloads[bsp.LumpPropHullVerts][i*12+j*4:], math.Float32bits(float32(i*100+j)))
		}
		binary.LittleEndian.PutUint32(payloads[bsp.LumpPropTris][i*8:], uint32(300+i))
		binary.LittleEndian.PutUint32(payloads[bsp.LumpPropTris][i*8+4:], uint32(15+i))
	}
	file := spatialFixture(t, bsp.ProfileReactiveDrop, payloads)
	collisions, err := file.PropCollisions()
	if err != nil || collisions[1].HullStart != 8 {
		t.Fatalf("collision scalars=%+v %v", collisions, err)
	}
	hulls, err := file.PropHulls()
	if err != nil || hulls[1].SurfaceProperty != -5 || hulls[1].Contents != 0x80000001 {
		t.Fatalf("hull fields=%+v %v", hulls, err)
	}
	vertices, err := file.PropHullVertices()
	if err != nil || vertices[1][2] != 102 {
		t.Fatalf("hull vertices=%+v %v", vertices, err)
	}
	triangles, err := file.PropHullTriangles()
	if err != nil || triangles[1].IndexStart != 301 {
		t.Fatalf("opaque-unit triangle scalars=%+v %v", triangles, err)
	}
	// A three-byte opaque blob does not establish an index-unit bound.
	blob, err := file.PropBlob()
	if err != nil || !bytes.Equal(blob, payloads[bsp.LumpPropBlob]) {
		t.Fatalf("byte blob=%x %v", blob, err)
	}
	before, _ := file.RawLump(bsp.LumpPropHulls)
	bad := append([]lumps.PropHull(nil), hulls...)
	bad[1].VertexCount = -1
	if err := file.SetPropHulls(bad); !errors.Is(err, bsp.ErrMalformed) {
		t.Fatalf("negative scalar span error=%v", err)
	}
	raw, _ := file.RawLump(bsp.LumpPropHulls)
	if !bytes.Equal(raw, before) {
		t.Fatal("bad hull scalar partially committed")
	}
	collisions[1].HullStart = 19
	hulls[1].SurfaceProperty = -71
	vertices[1][2] = -7.25
	triangles[1].IndexCount = 27
	blob[1] = 0x87
	if err := file.SetPropCollisions(collisions); err != nil {
		t.Fatal(err)
	}
	if err := file.SetPropHulls(hulls); err != nil {
		t.Fatal(err)
	}
	if err := file.SetPropHullVertices(vertices); err != nil {
		t.Fatal(err)
	}
	if err := file.SetPropHullTriangles(triangles); err != nil {
		t.Fatal(err)
	}
	if err := file.SetPropBlob(blob); err != nil {
		t.Fatal(err)
	}
	binary.LittleEndian.PutUint32(payloads[bsp.LumpPropCollision][12:], 19)
	binary.LittleEndian.PutUint32(payloads[bsp.LumpPropHulls][24:], 0xffffffb9)
	binary.LittleEndian.PutUint32(payloads[bsp.LumpPropHullVerts][20:], math.Float32bits(-7.25))
	binary.LittleEndian.PutUint32(payloads[bsp.LumpPropTris][12:], 27)
	payloads[bsp.LumpPropBlob][1] = 0x87
	next := rereadOrdinary(t, file)
	for id, expected := range payloads {
		raw, err := next.RawLump(id)
		if err != nil || !bytes.Equal(raw, expected) {
			t.Fatalf("RD%d committed wire bytes changed: %v", id, err)
		}
	}
	opaque := spatialFixture(t, bsp.ProfileCSGO, payloads)
	if _, err := opaque.PropCollisions(); !errors.Is(err, bsp.ErrUnsupportedLump) {
		t.Fatalf("historical numeric alias admitted: %v", err)
	}
	raw, _ = rereadOrdinary(t, opaque).RawLump(bsp.LumpPropCollision)
	if !bytes.Equal(raw, payloads[bsp.LumpPropCollision]) {
		t.Fatal("unsupported alias was not preserved")
	}
}

func TestTypedDisplacementEditPreservesAllOpaqueNumericIDs(t *testing.T) {
	for _, profile := range []bsp.Profile{bsp.ProfileSDK2013, bsp.ProfileCSGO, bsp.ProfileL4D2, bsp.ProfileReactiveDrop} {
		t.Run(string(profile), func(t *testing.T) {
			payloads := make(map[bsp.LumpID][]byte, 64)
			for id := bsp.LumpID(0); id < 64; id++ {
				payloads[id] = []byte{byte(id), 0xa5, 0x5a, byte(63 - id)}
			}
			payloads[bsp.LumpGame] = []byte{0, 0, 0, 0, 0xc1, 0xc2, 0xc3, 0xc4}
			geometry := displacementPayloads()
			for id, raw := range geometry {
				payloads[id] = raw
			}
			file := spatialFixture(t, profile, payloads)
			descriptors := make(map[bsp.LumpID]bsp.HeaderLump, 64)
			for id := bsp.LumpID(0); id < 64; id++ {
				desc, _ := file.Descriptor(id)
				if _, known := geometry[id]; !known && id != bsp.LumpGame {
					desc.Version = 77
					desc.FourCC = [4]byte{0x12, 0x34, byte(id), 0x56}
					if err := file.SetRawLump(id, desc, payloads[id]); err != nil {
						t.Fatal(err)
					}
				}
				descriptors[id], _ = file.Descriptor(id)
			}
			records, err := file.Displacements()
			if err != nil {
				t.Fatal(err)
			}
			records[1].AllowedVerts[9] = 0x76543210
			if err := file.SetDisplacements(records); err != nil {
				t.Fatal(err)
			}
			binary.LittleEndian.PutUint32(payloads[bsp.LumpDispInfo][176+172:], 0x76543210)
			next := rereadOrdinary(t, file)
			for id := bsp.LumpID(0); id < 64; id++ {
				raw, err := next.RawLump(id)
				desc, _ := next.Descriptor(id)
				expected := descriptors[id]
				if err != nil || !bytes.Equal(raw, payloads[id]) || desc.Version != expected.Version || desc.FourCC != expected.FourCC || desc.Length != int32(len(payloads[id])) {
					t.Fatalf("lump%d preservation raw=%x descriptor=%+v error=%v", id, raw, desc, err)
				}
			}
		})
	}
}

func TestZeroSampleOwnerDoesNotChangeSharedFollowingStream(t *testing.T) {
	payloads := displacementPayloads()
	binary.LittleEndian.PutUint32(payloads[bsp.LumpFaces][36:], 0xffffffff)
	binary.LittleEndian.PutUint32(payloads[bsp.LumpDispInfo][176+44:], 0)
	file := spatialFixture(t, bsp.ProfileCSGO, payloads)
	before := file.Header()
	original, _ := file.RawLump(bsp.LumpDispLightmapSamplePositions)
	if err := file.SetDisplacementSamples(0, nil); err != nil {
		t.Fatal(err)
	}
	after, _ := file.RawLump(bsp.LumpDispLightmapSamplePositions)
	following, err := file.DisplacementSamples(1)
	if err != nil || len(following) != 3 || following[2].Index != 510 || before != file.Header() || !bytes.Equal(original, after) {
		t.Fatalf("zero owner damaged shared populated stream: %+v %v", following, err)
	}
}

func TestDisplacementOwnerEditCannotOrphanMultiblend(t *testing.T) {
	payloads := displacementPayloads()
	binary.LittleEndian.PutUint32(payloads[bsp.LumpDispInfo][24:], lumps.DisplacementFlagMagic|lumps.DisplacementFlagHasMultiblend)
	payloads[bsp.LumpDispMultiBlend] = make([]byte, 289*80)
	file := spatialFixture(t, bsp.ProfileCSGO, payloads)
	records, err := file.Displacements()
	if err != nil {
		t.Fatal(err)
	}
	before, _ := file.RawLump(bsp.LumpDispInfo)
	records[0].Flags = lumps.DisplacementFlagMagic
	if err := file.SetDisplacements(records); !errors.Is(err, bsp.ErrMalformed) {
		t.Fatalf("orphaned multiblend error=%v", err)
	}
	after, _ := file.RawLump(bsp.LumpDispInfo)
	values, err := file.DisplacementMultiblends()
	if err != nil || len(values) != 289 || !bytes.Equal(before, after) {
		t.Fatalf("invalid owner edit committed: %v", err)
	}
}

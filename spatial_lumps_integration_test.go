package bsp_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"testing"

	"github.com/galaco/bsp"
)

func spatialFixture(t *testing.T, profile bsp.Profile, payloads map[bsp.LumpId][]byte) *bsp.Bsp {
	t.Helper()
	raw := make([]byte, bsp.HeaderSize)
	copy(raw, "VBSP")
	version := uint32(21)
	if profile == bsp.ProfileSDK2013 {
		version = 20
	}
	binary.LittleEndian.PutUint32(raw[4:], version)
	for id := bsp.LumpId(0); id < 64; id++ {
		payload := payloads[id]
		if len(payload) == 0 {
			continue
		}
		entry := raw[8+int(id)*16 : 8+int(id)*16+16]
		offset := uint32(len(raw))
		if profile == bsp.ProfileL4D2 {
			binary.LittleEndian.PutUint32(entry[4:], offset)
			binary.LittleEndian.PutUint32(entry[8:], uint32(len(payload)))
		} else {
			binary.LittleEndian.PutUint32(entry, offset)
			binary.LittleEndian.PutUint32(entry[4:], uint32(len(payload)))
		}
		raw = append(raw, payload...)
	}
	file, err := bsp.ReadFromStream(bytes.NewReader(raw), bsp.ReadOptions{Profile: profile})
	if err != nil {
		t.Fatal(err)
	}
	return file
}

func spatialPayloads() map[bsp.LumpId][]byte {
	overlays := make([]byte, 704)
	for i := 0; i < 2; i++ {
		wire := overlays[i*352 : (i+1)*352]
		binary.LittleEndian.PutUint32(wire, uint32(41+i))
		binary.LittleEndian.PutUint16(wire[6:], uint16(3<<14|2))
		for j := 0; j < 64; j++ {
			binary.LittleEndian.PutUint32(wire[8+j*4:], 0xffffffff)
		}
		binary.LittleEndian.PutUint32(wire[8:], 0)
		binary.LittleEndian.PutUint32(wire[12:], 1)
		for offset := 264; offset < 352; offset += 4 {
			binary.LittleEndian.PutUint32(wire[offset:], math.Float32bits(float32(offset+i*1000)))
		}
	}
	cubes := make([]byte, 32)
	for i, coordinates := range [][3]int32{{-12345, 16777217, 7}, {-2147483648, 2147483647, 88}} {
		for j, coordinate := range coordinates {
			binary.LittleEndian.PutUint32(cubes[i*16+j*4:], uint32(coordinate))
		}
		cubes[i*16+12] = byte(i * 5)
		copy(cubes[i*16+13:], []byte{0xa1, 0xb2, 0xc3})
	}
	primitives := make([]byte, 20)
	for i, fields := range [][4]uint16{{0, 3, 1, 3}, {3, 3, 4, 3}} {
		primitives[i*10], primitives[i*10+1] = byte(i), byte(0xa1+i)
		for j, value := range fields {
			binary.LittleEndian.PutUint16(primitives[i*10+2+j*2:], value)
		}
	}
	indices := make([]byte, 12)
	for i := 0; i < 6; i++ {
		binary.LittleEndian.PutUint16(indices[i*2:], uint16(i+1))
	}
	return map[bsp.LumpId][]byte{bsp.LumpOverlays: overlays, bsp.LumpCubemaps: cubes, bsp.LumpPrimitives: primitives, bsp.LumpPrimVerts: make([]byte, 7*12), bsp.LumpPrimIndices: indices, bsp.LumpTexInfo: make([]byte, 72), bsp.LumpFaces: make([]byte, 2*56)}
}

func TestSpatialPublicSnapshotsAndExactCommittedBytes(t *testing.T) {
	for _, profile := range []bsp.Profile{bsp.ProfileSDK2013, bsp.ProfileCSGO, bsp.ProfileL4D2} {
		t.Run(string(profile), func(t *testing.T) {
			payloads := spatialPayloads()
			file := spatialFixture(t, profile, payloads)
			overlays, err := file.Overlays()
			if err != nil {
				t.Fatal(err)
			}
			if len(overlays) != 2 || overlays[1].ID != 42 || overlays[1].UVPoints[3][2] != 1324 || overlays[1].BasisNormal[2] != 1348 || overlays[0].FaceCount() != 2 || overlays[0].RenderOrder() != 3 {
				t.Fatalf("complete overlay records=%+v", overlays)
			}
			overlays[1].UVPoints[3][2] = 7.25
			overlays[1].BasisNormal[2] = -0.5
			if err := overlays[1].SetRenderOrder(1); err != nil {
				t.Fatal(err)
			}
			fresh, err := file.Overlays()
			if err != nil || fresh[1].UVPoints[3][2] != 1324 {
				t.Fatalf("detached overlay mutation leaked: %v", err)
			}
			if err := file.SetOverlays(overlays); err != nil {
				t.Fatal(err)
			}
			binary.LittleEndian.PutUint16(payloads[bsp.LumpOverlays][352+6:], 1<<14|2)
			binary.LittleEndian.PutUint32(payloads[bsp.LumpOverlays][352+324:], math.Float32bits(7.25))
			binary.LittleEndian.PutUint32(payloads[bsp.LumpOverlays][352+348:], math.Float32bits(-0.5))
			cubes, err := file.Cubemaps()
			if err != nil {
				t.Fatal(err)
			}
			if len(cubes) != 2 || cubes[0].Origin != [3]int32{-12345, 16777217, 7} || cubes[1].Origin != [3]int32{-2147483648, 2147483647, 88} || cubes[1].Size != 5 || cubes[1].Padding != [3]byte{0xa1, 0xb2, 0xc3} {
				t.Fatalf("integer cubemaps=%+v", cubes)
			}
			cubes[1].Origin[2] = 16777219
			if err := file.SetCubemaps(cubes); err != nil {
				t.Fatal(err)
			}
			binary.LittleEndian.PutUint32(payloads[bsp.LumpCubemaps][24:], 16777219)
			primitives, err := file.Primitives()
			if err != nil {
				t.Fatal(err)
			}
			if len(primitives) != 2 || primitives[1].Padding != 0xa2 || primitives[1].FirstIndex != 3 || primitives[1].FirstVert != 4 || primitives[1].VertCount != 3 {
				t.Fatalf("aligned primitives=%+v", primitives)
			}
			primitives[1].FirstIndex, primitives[1].IndexCount = 4, 2
			if err := file.SetPrimitives(primitives); err != nil {
				t.Fatal(err)
			}
			binary.LittleEndian.PutUint16(payloads[bsp.LumpPrimitives][12:], 4)
			binary.LittleEndian.PutUint16(payloads[bsp.LumpPrimitives][14:], 2)
			next := rereadOrdinary(t, file)
			for _, id := range []bsp.LumpId{bsp.LumpOverlays, bsp.LumpCubemaps, bsp.LumpPrimitives} {
				got, err := next.RawLump(id)
				if err != nil || !bytes.Equal(got, payloads[id]) {
					t.Fatalf("lump%d exact committed bytes lost: %v", id, err)
				}
			}
			gotOverlays, err := next.Overlays()
			if err != nil || gotOverlays[1].UVPoints[3][2] != 7.25 || gotOverlays[1].Faces[63] != -1 || gotOverlays[1].RenderOrder() != 1 {
				t.Fatalf("overlay reread=%+v, %v", gotOverlays, err)
			}
			gotCubes, err := next.Cubemaps()
			if err != nil || gotCubes[1] != cubes[1] {
				t.Fatalf("cubemap reread=%+v, %v", gotCubes, err)
			}
			gotPrimitives, err := next.Primitives()
			if err != nil || gotPrimitives[1] != primitives[1] {
				t.Fatalf("primitive reread=%+v, %v", gotPrimitives, err)
			}
		})
	}
}

func TestSpatialSetterFailuresAreAtomic(t *testing.T) {
	file := spatialFixture(t, bsp.ProfileCSGO, spatialPayloads())
	overlays, err := file.Overlays()
	if err != nil {
		t.Fatal(err)
	}
	primitives, err := file.Primitives()
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name   string
		id     bsp.LumpId
		change func() error
	}{
		{"count65", bsp.LumpOverlays, func() error {
			bad := append(overlays[:0:0], overlays...)
			bad[1].FaceCountAndRenderOrder = 65
			return file.SetOverlays(bad)
		}},
		{"negative texinfo", bsp.LumpOverlays, func() error {
			bad := append(overlays[:0:0], overlays...)
			bad[1].TexInfo = -1
			return file.SetOverlays(bad)
		}},
		{"texinfo bound", bsp.LumpOverlays, func() error {
			bad := append(overlays[:0:0], overlays...)
			bad[1].TexInfo = 1
			return file.SetOverlays(bad)
		}},
		{"active face bound", bsp.LumpOverlays, func() error {
			bad := append(overlays[:0:0], overlays...)
			bad[1].Faces[0] = 2
			return file.SetOverlays(bad)
		}},
		{"index span", bsp.LumpPrimitives, func() error {
			bad := append(primitives[:0:0], primitives...)
			bad[1].IndexCount = 4
			return file.SetPrimitives(bad)
		}},
		{"vertex span", bsp.LumpPrimitives, func() error {
			bad := append(primitives[:0:0], primitives...)
			bad[1].VertCount = 4
			return file.SetPrimitives(bad)
		}},
		{"absolute index range", bsp.LumpPrimitives, func() error {
			bad := append(primitives[:0:0], primitives...)
			bad[1].FirstVert = 5
			bad[1].VertCount = 2
			return file.SetPrimitives(bad)
		}},
		{"unsupported type", bsp.LumpPrimitives, func() error {
			bad := append(primitives[:0:0], primitives...)
			bad[1].Type = 2
			err := file.SetPrimitives(bad)
			if !errors.Is(err, bsp.ErrUnsupportedLump) {
				t.Fatalf("unknown type error=%v", err)
			}
			return err
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			before, _ := file.RawLump(tc.id)
			descriptor, _ := file.Descriptor(tc.id)
			if err := tc.change(); err == nil {
				t.Fatal("invalid setter accepted")
			}
			after, _ := file.RawLump(tc.id)
			got, _ := file.Descriptor(tc.id)
			if !bytes.Equal(before, after) || got != descriptor {
				t.Fatal("failed setter changed bytes/descriptor")
			}
		})
	}
}

func TestPrimitiveFaceLocalIndicesAndOpaqueUnknownTypes(t *testing.T) {
	payloads := spatialPayloads()
	payloads[bsp.LumpPrimitives] = payloads[bsp.LumpPrimitives][:10]
	binary.LittleEndian.PutUint16(payloads[bsp.LumpPrimitives][6:], 0)
	binary.LittleEndian.PutUint16(payloads[bsp.LumpPrimitives][8:], 0)
	payloads[bsp.LumpPrimVerts] = nil
	payloads[bsp.LumpPrimIndices] = []byte{0, 0, 1, 0, 2, 0}
	binary.LittleEndian.PutUint16(payloads[bsp.LumpFaces][8:], 3)
	binary.LittleEndian.PutUint16(payloads[bsp.LumpFaces][48:], 0x8001)
	file := spatialFixture(t, bsp.ProfileCSGO, payloads)
	got, err := file.Primitives()
	if err != nil || len(got) != 1 || got[0].VertCount != 0 || got[0].IndexCount != 3 {
		t.Fatalf("valid face-local primitive=%+v, %v", got, err)
	}
	descriptor, _ := file.Descriptor(bsp.LumpPrimIndices)
	if err := file.SetRawLump(bsp.LumpPrimIndices, descriptor, []byte{0, 0, 1, 0, 3, 0}); err != nil {
		t.Fatal(err)
	}
	if _, err := file.Primitives(); err == nil {
		t.Fatal("face-local index beyond winding accepted")
	}
	unknown := append([]byte(nil), payloads[bsp.LumpPrimitives]...)
	unknown[0] = 99
	file = spatialFixture(t, bsp.ProfileCSGO, map[bsp.LumpId][]byte{bsp.LumpPrimitives: unknown})
	if _, err := file.Primitives(); !errors.Is(err, bsp.ErrUnsupportedLump) {
		t.Fatalf("unknown typed primitive=%v", err)
	}
	raw, err := rereadOrdinary(t, file).RawLump(bsp.LumpPrimitives)
	if err != nil || !bytes.Equal(raw, unknown) {
		t.Fatalf("opaque primitive not preserved: %v", err)
	}
}

func TestPublicFixedRecordBoundariesRetainOpaqueBytes(t *testing.T) {
	inventory := map[bsp.LumpId]int{1: 20, 2: 32, 3: 12, 5: 32, 6: 72, 7: 56, 8: 4, 11: 2, 12: 4, 13: 4, 14: 48, 16: 2, 17: 2, 18: 12, 20: 8, 21: 12, 26: 176, 27: 56, 30: 12, 31: 2, 33: 20, 36: 12, 37: 10, 38: 12, 39: 2, 41: 12, 42: 16, 44: 4, 45: 352, 46: 2, 47: 2, 48: 2, 53: 4, 58: 56, 60: 8}
	for id, width := range inventory {
		for _, length := range []int{width - 1, width + 1} {
			raw := make([]byte, length)
			file := spatialFixture(t, bsp.ProfileCSGO, map[bsp.LumpId][]byte{id: raw})
			if _, err := file.Lump(id); !errors.Is(err, bsp.ErrMalformed) {
				t.Fatalf("lump%d length%d error=%v", id, length, err)
			}
			preserved, err := rereadOrdinary(t, file).RawLump(id)
			if err != nil || !bytes.Equal(preserved, raw) {
				t.Fatalf("lump%d malformed opaque export altered: %v", id, err)
			}
		}
	}
}

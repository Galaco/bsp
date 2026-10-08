package codec

import (
	"bytes"
	"encoding/binary"
	"math"
	"reflect"
	"testing"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

func coverageRDContext() format.Context {
	return format.Context{Profile: "reactivedrop", BSPVersion: 21, ByteOrder: binary.LittleEndian}
}

// Independently authored bytes at the specified 1120-byte offsets. Record 0
// has inactive data through face 255; record 1 uses every face slot.
func coverageWaterFixture() []byte {
	raw := make([]byte, 2240)
	for i := 0; i < 2; i++ {
		wire := raw[i*1120 : (i+1)*1120]
		binary.LittleEndian.PutUint32(wire[0:4], 0x80000007+uint32(i))
		binary.LittleEndian.PutUint16(wire[4:6], 0xfffe+uint16(i))
		binary.LittleEndian.PutUint16(wire[6:8], []uint16{0x8002, 0xc100}[i])
		for j := 0; j < 256; j++ {
			binary.LittleEndian.PutUint32(wire[8+j*4:12+j*4], 0xffff0000+uint32(i*256+j))
		}
		for j := 0; j < 22; j++ {
			binary.LittleEndian.PutUint32(wire[1032+j*4:1036+j*4], 0x3f000001+uint32(i*32+j))
		}
		binary.LittleEndian.PutUint32(wire[1032:1036], 0x80000000)
		binary.LittleEndian.PutUint32(wire[1040:1044], 0x7f800000)
		binary.LittleEndian.PutUint32(wire[1052:1056], 0x7fc01234+uint32(i))
		binary.LittleEndian.PutUint32(wire[1112:1116], 0xffc02345+uint32(i))
	}
	return raw
}

func TestCoverageWaterIndependentWire(t *testing.T) {
	contexts := append(spatialContexts(), coverageRDContext())
	for _, context := range contexts {
		raw := coverageWaterFixture()
		records, err := DecodeWaterOverlays(raw, context)
		if err != nil || len(records) != 2 {
			t.Fatalf("%+v: len=%d err=%v", context, len(records), err)
		}
		for i := range records {
			record := &records[i]
			if record.ID != -2147483641+int32(i) || record.TexInfo != -2+int16(i) || record.FaceCount() != []uint16{2, 256}[i] || record.RenderOrder() != uint16(i+2) {
				t.Fatalf("record %d: incorrect signed/packed header: %+v", i, record)
			}
			for j, face := range record.Faces {
				if face != -65536+int32(i*256+j) {
					t.Fatalf("record %d face %d: %d", i, j, face)
				}
			}
			floats := []float32{record.U[0], record.U[1], record.V[0], record.V[1]}
			for _, point := range record.UVPoints {
				floats = append(floats, point[:]...)
			}
			floats = append(floats, record.Origin[:]...)
			floats = append(floats, record.BasisNormal[:]...)
			for j, value := range floats {
				want := uint32(0x3f000001 + i*32 + j)
				switch j {
				case 0:
					want = 0x80000000
				case 2:
					want = 0x7f800000
				case 5:
					want = 0x7fc01234 + uint32(i)
				case 20:
					want = 0xffc02345 + uint32(i)
				}
				if bits := math.Float32bits(value); bits != want {
					t.Fatalf("record %d geometry component %d: %08x want %08x", i, j, bits, want)
				}
			}
		}
		encoded, err := EncodeWaterOverlays(records, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("water roundtrip changed bytes: %v", err)
		}
		records[1].Faces[255] = 17
		records[1].UVPoints[3][2] = -2
		records[1].BasisNormal[2] = 4
		binary.LittleEndian.PutUint32(raw[1120+1028:], 17)
		binary.LittleEndian.PutUint32(raw[1120+1092:], 0xc0000000)
		binary.LittleEndian.PutUint32(raw[1120+1116:], 0x40800000)
		encoded, err = EncodeWaterOverlays(records, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("water later-field edit changed unrelated bytes: %v", err)
		}
	}
}

func TestCoverageWaterCountOrder(t *testing.T) {
	record := lumps.WaterOverlay{FaceCountAndRenderOrder: 0xc100}
	record.Faces[255] = -42
	if err := record.SetFaceCount(1); err != nil || record.FaceCount() != 1 || record.RenderOrder() != 3 || record.Faces[255] != -42 {
		t.Fatalf("count setter changed order/inactive data: %+v %v", record, err)
	}
	if err := record.SetRenderOrder(2); err != nil || record.FaceCount() != 1 || record.RenderOrder() != 2 {
		t.Fatalf("order setter changed count: %+v %v", record, err)
	}
	before := record
	for _, count := range []uint16{257, 0x3fff, 0x4000, 0xffff} {
		if err := record.SetFaceCount(count); err == nil || record != before {
			t.Fatalf("invalid count %d accepted or mutated record", count)
		}
	}
	for _, order := range []uint16{4, 0x8000, 0xffff} {
		if err := record.SetRenderOrder(order); err == nil || record != before {
			t.Fatalf("invalid order %d accepted or mutated record", order)
		}
	}
	if err := record.SetFaceCount(256); err != nil || record.FaceCount() != 256 || record.RenderOrder() != 2 {
		t.Fatalf("maximum count rejected: %v", err)
	}
	for _, count := range []uint16{257, 0x3fff} {
		raw := coverageWaterFixture()
		binary.LittleEndian.PutUint16(raw[1120+6:], 0xc000|count)
		if _, err := DecodeWaterOverlays(raw, coverageRDContext()); err == nil {
			t.Fatalf("decode accepted count %d", count)
		}
		record.FaceCountAndRenderOrder = 0xc000 | count
		if _, err := EncodeWaterOverlays([]lumps.WaterOverlay{record}, coverageRDContext()); err == nil {
			t.Fatalf("encode accepted count %d", count)
		}
	}
}

func TestCoverageSystemLevelsIndependentWire(t *testing.T) {
	raw := []byte{0xff, 3, 0x80, 7, 2, 0xfe, 9, 1}
	want := []lumps.OverlaySystemLevel{
		{MinCPULevel: 255, MaxCPULevel: 3, MinGPULevel: 128, MaxGPULevel: 7},
		{MinCPULevel: 2, MaxCPULevel: 254, MinGPULevel: 9, MaxGPULevel: 1},
	}
	for _, context := range []format.Context{{Profile: "csgo", BSPVersion: 21}, {Profile: "l4d2", BSPVersion: 21}, coverageRDContext()} {
		records, err := DecodeOverlaySystemLevels(raw, context)
		if err != nil || !reflect.DeepEqual(records, want) {
			t.Fatalf("asymmetric levels lost: %+v %v", records, err)
		}
		encoded, err := EncodeOverlaySystemLevels(records, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("levels changed bytes: %v", err)
		}
		records[1].MinGPULevel = 42
		encoded, err = EncodeOverlaySystemLevels(records, context)
		if err != nil || !bytes.Equal(encoded, []byte{0xff, 3, 0x80, 7, 2, 0xfe, 42, 1}) {
			t.Fatalf("levels edit changed wrong byte: %v", err)
		}
	}
}

func coverageMultiblendFixture() []byte {
	raw := make([]byte, 160)
	for i := 0; i < 2; i++ {
		for j := 0; j < 20; j++ {
			binary.LittleEndian.PutUint32(raw[i*80+j*4:], 0x3f000011+uint32(i*32+j))
		}
		binary.LittleEndian.PutUint32(raw[i*80:], 0x80000000)
		binary.LittleEndian.PutUint32(raw[i*80+28:], 0xff800000)
		binary.LittleEndian.PutUint32(raw[i*80+76:], 0x7fc12345+uint32(i))
	}
	return raw
}

func TestCoverageMultiblendIndependentWire(t *testing.T) {
	for _, context := range []format.Context{{Profile: "csgo", BSPVersion: 21}, coverageRDContext()} {
		raw := coverageMultiblendFixture()
		records, err := DecodeDisplacementMultiblends(raw, context)
		if err != nil || len(records) != 2 {
			t.Fatalf("multiblend decode: len=%d err=%v", len(records), err)
		}
		for i := range records {
			record := &records[i]
			floats := append([]float32(nil), record.Weights[:]...)
			floats = append(floats, record.Alpha[:]...)
			for _, color := range record.Colors {
				floats = append(floats, color[:]...)
			}
			for j, value := range floats {
				want := uint32(0x3f000011 + i*32 + j)
				switch j {
				case 0:
					want = 0x80000000
				case 7:
					want = 0xff800000
				case 19:
					want = 0x7fc12345 + uint32(i)
				}
				if bits := math.Float32bits(value); bits != want {
					t.Fatalf("record %d component %d: %08x want %08x", i, j, bits, want)
				}
			}
		}
		encoded, err := EncodeDisplacementMultiblends(records, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("multiblend changed bytes: %v", err)
		}
		records[1].Colors[3][2] = -2
		binary.LittleEndian.PutUint32(raw[156:], 0xc0000000)
		encoded, err = EncodeDisplacementMultiblends(records, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("last color edit changed unrelated bytes: %v", err)
		}
	}
}

func TestCoveragePropIndependentWire(t *testing.T) {
	context := coverageRDContext()
	collisionRaw := []byte{2, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 0xff, 0x7f}
	collisionWant := []lumps.PropCollision{{HullCount: 2, HullStart: 7}, {HullCount: 0, HullStart: math.MaxInt32}}
	collisions, err := DecodePropCollisions(collisionRaw, context)
	if err != nil || !reflect.DeepEqual(collisions, collisionWant) {
		t.Fatalf("prop collisions: %+v %v", collisions, err)
	}
	if encoded, err := EncodePropCollisions(collisions, context); err != nil || !bytes.Equal(encoded, collisionRaw) {
		t.Fatalf("collision bytes: %v", err)
	}
	collisions[1].HullStart = 0x12345678
	copy(collisionRaw[12:], []byte{0x78, 0x56, 0x34, 0x12})
	if encoded, err := EncodePropCollisions(collisions, context); err != nil || !bytes.Equal(encoded, collisionRaw) {
		t.Fatalf("collision second-record edit: %v", err)
	}

	hullRaw := []byte{
		3, 0, 0, 0, 17, 0, 0, 0, 0xfe, 0xff, 0xff, 0xff, 1, 0, 0, 0x80,
		0, 0, 0, 0, 0xff, 0xff, 0xff, 0x7f, 0, 0, 0, 0x80, 0xff, 0xff, 0xff, 0xff,
	}
	hullWant := []lumps.PropHull{
		{VertexCount: 3, VertexStart: 17, SurfaceProperty: -2, Contents: 0x80000001},
		{VertexCount: 0, VertexStart: math.MaxInt32, SurfaceProperty: math.MinInt32, Contents: math.MaxUint32},
	}
	hulls, err := DecodePropHulls(hullRaw, context)
	if err != nil || !reflect.DeepEqual(hulls, hullWant) {
		t.Fatalf("signed hull fields/unsigned contents: %+v %v", hulls, err)
	}
	if encoded, err := EncodePropHulls(hulls, context); err != nil || !bytes.Equal(encoded, hullRaw) {
		t.Fatalf("hull bytes: %v", err)
	}
	hulls[1].SurfaceProperty = -7
	copy(hullRaw[24:], []byte{0xf9, 0xff, 0xff, 0xff})
	if encoded, err := EncodePropHulls(hulls, context); err != nil || !bytes.Equal(encoded, hullRaw) {
		t.Fatalf("hull second-record edit: %v", err)
	}

	vertexRaw := []byte{0, 0, 0, 0x80, 0, 0, 0x80, 0x7f, 0x34, 0x12, 0xc0, 0x7f, 0, 0, 0, 0x40, 0, 0, 0, 0xc0, 0x45, 0x23, 0xc0, 0xff}
	vertices, err := DecodePropHullVertices(vertexRaw, context)
	if err != nil || len(vertices) != 2 {
		t.Fatalf("hull vertices: %+v %v", vertices, err)
	}
	for i := range vertices {
		for j, value := range vertices[i] {
			if math.Float32bits(value) != binary.LittleEndian.Uint32(vertexRaw[i*12+j*4:]) {
				t.Fatalf("vertex %d coordinate %d bits changed", i, j)
			}
		}
	}
	if encoded, err := EncodePropHullVertices(vertices, context); err != nil || !bytes.Equal(encoded, vertexRaw) {
		t.Fatalf("vertex bytes: %v", err)
	}
	vertices[1][2] = 4
	copy(vertexRaw[20:], []byte{0, 0, 0x80, 0x40})
	if encoded, err := EncodePropHullVertices(vertices, context); err != nil || !bytes.Equal(encoded, vertexRaw) {
		t.Fatalf("vertex second-record edit: %v", err)
	}

	triangleRaw := []byte{31, 0, 0, 0, 6, 0, 0, 0, 0xff, 0xff, 0xff, 0x7f, 0, 0, 0, 0}
	triangleWant := []lumps.PropHullTriangles{{IndexStart: 31, IndexCount: 6}, {IndexStart: math.MaxInt32, IndexCount: 0}}
	triangles, err := DecodePropHullTriangles(triangleRaw, context)
	if err != nil || !reflect.DeepEqual(triangles, triangleWant) {
		t.Fatalf("triangle spans: %+v %v", triangles, err)
	}
	if encoded, err := EncodePropHullTriangles(triangles, context); err != nil || !bytes.Equal(encoded, triangleRaw) {
		t.Fatalf("triangle bytes: %v", err)
	}
	triangles[1].IndexStart = 0x12345678
	copy(triangleRaw[8:], []byte{0x78, 0x56, 0x34, 0x12})
	if encoded, err := EncodePropHullTriangles(triangles, context); err != nil || !bytes.Equal(encoded, triangleRaw) {
		t.Fatalf("triangle second-record edit: %v", err)
	}
}

func TestCoveragePropInvalidSelfSpans(t *testing.T) {
	context := coverageRDContext()
	for _, span := range [][2]int32{{-1, 0}, {0, -1}, {math.MinInt32, 1}, {math.MaxInt32, 1}, {1, math.MaxInt32}} {
		start, count := span[0], span[1]
		for _, tc := range []struct {
			name                     string
			width                    int
			startOffset, countOffset int
			decode                   func([]byte, format.Context) error
			encode                   func() error
		}{
			{"collision", 8, 4, 0, func(raw []byte, c format.Context) error { _, err := DecodePropCollisions(raw, c); return err }, func() error {
				_, err := EncodePropCollisions([]lumps.PropCollision{{HullCount: count, HullStart: start}}, context)
				return err
			}},
			{"hull", 16, 4, 0, func(raw []byte, c format.Context) error { _, err := DecodePropHulls(raw, c); return err }, func() error {
				_, err := EncodePropHulls([]lumps.PropHull{{VertexCount: count, VertexStart: start}}, context)
				return err
			}},
			{"triangles", 8, 0, 4, func(raw []byte, c format.Context) error { _, err := DecodePropHullTriangles(raw, c); return err }, func() error {
				_, err := EncodePropHullTriangles([]lumps.PropHullTriangles{{IndexStart: start, IndexCount: count}}, context)
				return err
			}},
		} {
			raw := make([]byte, tc.width*2)
			binary.LittleEndian.PutUint32(raw[tc.width+tc.startOffset:], uint32(start))
			binary.LittleEndian.PutUint32(raw[tc.width+tc.countOffset:], uint32(count))
			if err := tc.decode(raw, context); err == nil {
				t.Fatalf("%s accepted malformed second span %v", tc.name, span)
			}
			if err := tc.encode(); err == nil {
				t.Fatalf("%s encoded malformed span %v", tc.name, span)
			}
		}
	}
}

func TestCoverageExactWidthsAndProfileGates(t *testing.T) {
	cases := []struct {
		name     string
		width    int
		profiles map[string]bool
		decode   func([]byte, format.Context) error
		encode   func(format.Context) error
	}{
		{"water", 1120, map[string]bool{"sdk2013": true, "csgo": true, "l4d2": true, "reactivedrop": true}, func(raw []byte, c format.Context) error { _, err := DecodeWaterOverlays(raw, c); return err }, func(c format.Context) error { _, err := EncodeWaterOverlays(nil, c); return err }},
		{"levels", 4, map[string]bool{"csgo": true, "l4d2": true, "reactivedrop": true}, func(raw []byte, c format.Context) error { _, err := DecodeOverlaySystemLevels(raw, c); return err }, func(c format.Context) error { _, err := EncodeOverlaySystemLevels(nil, c); return err }},
		{"multiblend", 80, map[string]bool{"csgo": true, "reactivedrop": true}, func(raw []byte, c format.Context) error { _, err := DecodeDisplacementMultiblends(raw, c); return err }, func(c format.Context) error { _, err := EncodeDisplacementMultiblends(nil, c); return err }},
		{"collision", 8, map[string]bool{"reactivedrop": true}, func(raw []byte, c format.Context) error { _, err := DecodePropCollisions(raw, c); return err }, func(c format.Context) error { _, err := EncodePropCollisions(nil, c); return err }},
		{"hull", 16, map[string]bool{"reactivedrop": true}, func(raw []byte, c format.Context) error { _, err := DecodePropHulls(raw, c); return err }, func(c format.Context) error { _, err := EncodePropHulls(nil, c); return err }},
		{"vertex", 12, map[string]bool{"reactivedrop": true}, func(raw []byte, c format.Context) error { _, err := DecodePropHullVertices(raw, c); return err }, func(c format.Context) error { _, err := EncodePropHullVertices(nil, c); return err }},
		{"triangles", 8, map[string]bool{"reactivedrop": true}, func(raw []byte, c format.Context) error { _, err := DecodePropHullTriangles(raw, c); return err }, func(c format.Context) error { _, err := EncodePropHullTriangles(nil, c); return err }},
	}
	contexts := append(spatialContexts(), coverageRDContext(), format.Context{Profile: "reactivedrop", BSPVersion: 19}, format.Context{Profile: "reactivedrop", BSPVersion: 20})
	for _, tc := range cases {
		for _, context := range contexts {
			accepted := tc.profiles[context.Profile]
			for _, size := range []int{0, tc.width, tc.width * 2} {
				if err := tc.decode(make([]byte, size), context); (err == nil) != accepted {
					t.Fatalf("%s context %+v size %d admission err=%v", tc.name, context, size, err)
				}
			}
			if err := tc.encode(context); (err == nil) != accepted {
				t.Fatalf("%s encode context %+v admission err=%v", tc.name, context, err)
			}
			if !accepted {
				continue
			}
			for size := 1; size < 2*tc.width; size++ {
				if size%tc.width != 0 {
					if err := tc.decode(make([]byte, size), context); err == nil {
						t.Fatalf("%s accepted incomplete length %d", tc.name, size)
					}
				}
			}
			for _, invalid := range []format.Context{
				{Profile: context.Profile, BSPVersion: context.BSPVersion, LumpVersion: 1},
				{Profile: context.Profile, BSPVersion: context.BSPVersion, LumpVersion: -1},
				{Profile: context.Profile, BSPVersion: context.BSPVersion, ByteOrder: binary.BigEndian},
			} {
				if err := tc.decode(nil, invalid); err == nil {
					t.Fatalf("%s decoded unsupported context %+v", tc.name, invalid)
				}
				if err := tc.encode(invalid); err == nil {
					t.Fatalf("%s encoded unsupported context %+v", tc.name, invalid)
				}
			}
		}
		for _, invalid := range []format.Context{{}, {Profile: "unknown", BSPVersion: 21}, {Profile: "sdk2013", BSPVersion: 21}, {Profile: "csgo", BSPVersion: 20}, {Profile: "l4d2", BSPVersion: 20}, {Profile: "reactivedrop", BSPVersion: 18}, {Profile: "reactivedrop", BSPVersion: 22}} {
			if err := tc.decode(nil, invalid); err == nil {
				t.Fatalf("%s decoded unverified context %+v", tc.name, invalid)
			}
			if err := tc.encode(invalid); err == nil {
				t.Fatalf("%s encoded unverified context %+v", tc.name, invalid)
			}
		}
	}
}

func TestCoverageAdaptersRetainOnFailureAndResetEmpty(t *testing.T) {
	context := coverageRDContext()
	water := &WaterOverlays{}
	levels := &OverlaySystemLevels{}
	multiblend := &DisplacementMultiblends{}
	collision := &PropCollisions{}
	hull := &PropHulls{}
	vertex := &PropHullVertices{}
	triangles := &PropHullTriangles{}
	cases := []struct {
		name    string
		adapter interface {
			SetContext(format.Context)
			Length() int
			Unmarshall([]byte) error
			Marshall() ([]byte, error)
		}
		raw    []byte
		mutate func()
		offset int
		bits   uint32
	}{
		{"water", water, coverageWaterFixture(), func() { water.GetData()[1].Faces[255] = 17 }, 2148, 17},
		{"levels", levels, []byte{1, 2, 3, 4, 5, 6, 7, 8}, func() { levels.GetData()[1].MinGPULevel = 42 }, 4, 0x082a0605},
		{"multiblend", multiblend, coverageMultiblendFixture(), func() { multiblend.GetData()[1].Colors[3][2] = -2 }, 156, 0xc0000000},
		{"collision", collision, []byte{1, 0, 0, 0, 2, 0, 0, 0, 3, 0, 0, 0, 4, 0, 0, 0}, func() { collision.GetData()[1].HullStart = 17 }, 12, 17},
		{"hull", hull, []byte{1, 0, 0, 0, 2, 0, 0, 0, 3, 0, 0, 0, 4, 0, 0, 0, 5, 0, 0, 0, 6, 0, 0, 0, 7, 0, 0, 0, 8, 0, 0, 0}, func() { hull.GetData()[1].Contents = 0x80000001 }, 28, 0x80000001},
		{"vertex", vertex, []byte{0, 0, 0x80, 0x3f, 0, 0, 0, 0x40, 0, 0, 0x40, 0x40, 0, 0, 0x80, 0x40, 0, 0, 0xa0, 0x40, 0, 0, 0xc0, 0x40}, func() { vertex.GetData()[1][2] = -2 }, 20, 0xc0000000},
		{"triangles", triangles, []byte{1, 0, 0, 0, 2, 0, 0, 0, 3, 0, 0, 0, 4, 0, 0, 0}, func() { triangles.GetData()[1].IndexStart = 17 }, 8, 17},
	}
	for _, tc := range cases {
		tc.adapter.SetContext(context)
		if err := tc.adapter.Unmarshall(tc.raw); err != nil {
			t.Fatalf("%s decode: %v", tc.name, err)
		}
		tc.mutate()
		binary.LittleEndian.PutUint32(tc.raw[tc.offset:], tc.bits)
		if err := tc.adapter.Unmarshall([]byte{1}); err == nil || tc.adapter.Length() != len(tc.raw) {
			t.Fatalf("%s bad decode accepted/metadata mutated", tc.name)
		}
		if encoded, err := tc.adapter.Marshall(); err != nil || !bytes.Equal(encoded, tc.raw) {
			t.Fatalf("%s adapter lost edited data on failure: %v", tc.name, err)
		}
		if err := tc.adapter.Unmarshall(nil); err != nil || tc.adapter.Length() != 0 {
			t.Fatalf("%s empty reset: %v", tc.name, err)
		}
		if encoded, err := tc.adapter.Marshall(); err != nil || len(encoded) != 0 {
			t.Fatalf("%s stale records after empty reset: %v", tc.name, err)
		}
	}
}

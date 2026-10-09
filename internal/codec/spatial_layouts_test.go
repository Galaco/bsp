package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

func spatialContexts() []format.Context {
	return []format.Context{
		{Profile: "sdk2013", BSPVersion: 19, ByteOrder: binary.LittleEndian},
		{Profile: "sdk2013", BSPVersion: 20},
		{Profile: "l4d2", BSPVersion: 21, ByteOrder: binary.LittleEndian},
		{Profile: "csgo", BSPVersion: 21, ByteOrder: binary.LittleEndian},
	}
}

// Fixtures are authored as bytes at documented offsets, never by encoding
// production records. Every float slot has its own independently known bits.
func spatialOverlayFixture() []byte {
	raw := make([]byte, 704)
	for i := 0; i < 2; i++ {
		wire := raw[i*352 : (i+1)*352]
		binary.LittleEndian.PutUint32(wire[0:4], 0x80000007+uint32(i))
		binary.LittleEndian.PutUint16(wire[4:6], 0xfffe+uint16(i))
		binary.LittleEndian.PutUint16(wire[6:8], 0xc002+uint16(i))
		for j := 0; j < 64; j++ {
			binary.LittleEndian.PutUint32(wire[8+j*4:12+j*4], 0xffff0000+uint32(i*64+j))
		}
		for j := 0; j < 22; j++ {
			binary.LittleEndian.PutUint32(wire[264+j*4:268+j*4], 0x3f000001+uint32(i*32+j))
		}
		// Negative zero, infinity and noncanonical NaN payloads must survive.
		binary.LittleEndian.PutUint32(wire[264:268], 0x80000000)
		binary.LittleEndian.PutUint32(wire[272:276], 0x7f800000)
		binary.LittleEndian.PutUint32(wire[284:288], 0x7fc01234+uint32(i))
		binary.LittleEndian.PutUint32(wire[344:348], 0xffc02345+uint32(i))
	}
	return raw
}

func TestSpatialOverlaysIndependentWire(t *testing.T) {
	for _, context := range spatialContexts() {
		raw := spatialOverlayFixture()
		records, err := DecodeOverlays(raw, context)
		if err != nil || len(records) != 2 {
			t.Fatalf("%+v: decode len=%d err=%v", context, len(records), err)
		}
		for i, record := range records {
			if record.ID != -2147483641+int32(i) || record.TexInfo != -2+int16(i) || record.FaceCount() != 2+uint16(i) || record.RenderOrder() != 3 {
				t.Fatalf("record %d: incorrect signed header or packed fields: %+v", i, record)
			}
			for j, face := range record.Faces {
				if face != -65536+int32(i*64+j) {
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
					t.Fatalf("record %d float %d: %08x want %08x", i, j, bits, want)
				}
			}
		}
		encoded, err := EncodeOverlays(records, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("overlay bytes changed: %v", err)
		}
		records[1].UVPoints[3][2] = -2
		records[1].BasisNormal[2] = 4
		binary.LittleEndian.PutUint32(raw[352+324:], 0xc0000000)
		binary.LittleEndian.PutUint32(raw[352+348:], 0x40800000)
		encoded, err = EncodeOverlays(records, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("later-field edits changed unrelated bytes: %v", err)
		}
	}
}

func spatialCubemapFixture() []byte {
	return []byte{
		0xc7, 0xcf, 0xff, 0xff, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x80, 0xff, 0xa1, 0xb2, 0xc3,
		0xff, 0xff, 0xff, 0x7f, 0xff, 0xff, 0xff, 0xff, 0x2a, 0x00, 0x00, 0x00, 0x00, 0xd4, 0xe5, 0xf6,
	}
}

func TestSpatialCubemapsIndependentWire(t *testing.T) {
	want := []lumps.Cubemap{
		{Origin: [3]int32{-12345, 16777217, -2147483648}, Size: 255, Padding: [3]byte{0xa1, 0xb2, 0xc3}},
		{Origin: [3]int32{2147483647, -1, 42}, Size: 0, Padding: [3]byte{0xd4, 0xe5, 0xf6}},
	}
	for _, context := range spatialContexts() {
		raw := spatialCubemapFixture()
		records, err := DecodeCubemaps(raw, context)
		if err != nil || !reflect.DeepEqual(records, want) {
			t.Fatalf("cubemap decode: %+v %v", records, err)
		}
		encoded, err := EncodeCubemaps(want, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("cubemap encode: %v", err)
		}
		records[1].Origin[2] = -7
		copy(raw[24:28], []byte{0xf9, 0xff, 0xff, 0xff})
		encoded, err = EncodeCubemaps(records, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("cubemap edit: %v", err)
		}
	}
}

func spatialPrimitiveFixture() []byte {
	return []byte{0, 0xa5, 0x34, 0x12, 0x03, 0x00, 0x78, 0x56, 0x04, 0x00,
		1, 0x5a, 0xcd, 0xab, 0x08, 0x00, 0xff, 0xff, 0x07, 0x00}
}

func TestSpatialPrimitivesIndependentWire(t *testing.T) {
	want := []lumps.Primitive{
		{Type: 0, Padding: 0xa5, FirstIndex: 0x1234, IndexCount: 3, FirstVert: 0x5678, VertCount: 4},
		{Type: 1, Padding: 0x5a, FirstIndex: 0xabcd, IndexCount: 8, FirstVert: 0xffff, VertCount: 7},
	}
	for _, context := range spatialContexts() {
		raw := spatialPrimitiveFixture()
		records, err := DecodePrimitives(raw, context)
		if err != nil || !reflect.DeepEqual(records, want) {
			t.Fatalf("primitive decode: %+v %v", records, err)
		}
		encoded, err := EncodePrimitives(want, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("primitive encode: %v", err)
		}
		records[1].VertCount = 0x2345
		raw[18], raw[19] = 0x45, 0x23
		encoded, err = EncodePrimitives(records, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("primitive edit: %v", err)
		}
	}
}

func TestSpatialExactSizesEmptyAndContexts(t *testing.T) {
	cases := []struct {
		name        string
		width       int
		decode      func([]byte, format.Context) error
		encodeEmpty func(format.Context) error
	}{
		{"overlay", 352, func(raw []byte, c format.Context) error { _, err := DecodeOverlays(raw, c); return err }, func(c format.Context) error { _, err := EncodeOverlays(nil, c); return err }},
		{"cubemap", 16, func(raw []byte, c format.Context) error { _, err := DecodeCubemaps(raw, c); return err }, func(c format.Context) error { _, err := EncodeCubemaps(nil, c); return err }},
		{"primitive", 10, func(raw []byte, c format.Context) error { _, err := DecodePrimitives(raw, c); return err }, func(c format.Context) error { _, err := EncodePrimitives(nil, c); return err }},
	}
	invalid := []format.Context{
		{}, {Profile: "unknown", BSPVersion: 20}, {Profile: "sdk2013", BSPVersion: 21},
		{Profile: "csgo", BSPVersion: 20}, {Profile: "l4d2", BSPVersion: 19},
		{Profile: "sdk2013", BSPVersion: 20, LumpVersion: 1},
		{Profile: "csgo", BSPVersion: 21, LumpVersion: -1},
		{Profile: "l4d2", BSPVersion: 21, LumpVersion: 1},
		{Profile: "sdk2013", BSPVersion: 20, ByteOrder: binary.BigEndian},
	}
	for _, tc := range cases {
		for _, context := range spatialContexts() {
			for _, size := range []int{1, tc.width - 1, tc.width + 1, 2*tc.width - 1, 2*tc.width + 1} {
				if err := tc.decode(make([]byte, size), context); err == nil {
					t.Fatalf("%s accepted length %d", tc.name, size)
				}
			}
			for _, size := range []int{0, tc.width, 2 * tc.width} {
				if err := tc.decode(make([]byte, size), context); err != nil {
					t.Fatalf("%s length %d: %v", tc.name, size, err)
				}
			}
			if err := tc.encodeEmpty(context); err != nil {
				t.Fatalf("%s empty encode: %v", tc.name, err)
			}
		}
		for _, context := range invalid {
			if err := tc.decode(nil, context); err == nil || errors.Is(err, ErrUnsupportedPrimitiveType) {
				t.Fatalf("%s wrong context error: %+v %v", tc.name, context, err)
			}
			if err := tc.encodeEmpty(context); err == nil || errors.Is(err, ErrUnsupportedPrimitiveType) {
				t.Fatalf("%s wrong encode context error: %+v %v", tc.name, context, err)
			}
		}
	}
}

func TestSpatialPackedEditsAreAtomic(t *testing.T) {
	record := lumps.Overlay{FaceCountAndRenderOrder: 0xc040}
	for _, count := range []uint16{65, 0x3fff, 0xffff} {
		before := record
		if err := record.SetFaceCount(count); err == nil || record != before {
			t.Fatalf("count %d not rejected atomically", count)
		}
	}
	for _, order := range []uint16{4, 0xffff} {
		before := record
		if err := record.SetRenderOrder(order); err == nil || record != before {
			t.Fatalf("order %d not rejected atomically", order)
		}
	}
	for _, count := range []uint16{0, 1, 64} {
		if err := record.SetFaceCount(count); err != nil || record.FaceCount() != count || record.RenderOrder() != 3 {
			t.Fatalf("valid count edit: %+v %v", record, err)
		}
	}
	for _, order := range []uint16{0, 1, 2, 3} {
		if err := record.SetRenderOrder(order); err != nil || record.FaceCount() != 64 || record.RenderOrder() != order {
			t.Fatalf("valid order edit: %+v %v", record, err)
		}
	}
	context := spatialContexts()[0]
	raw := spatialOverlayFixture()
	binary.LittleEndian.PutUint16(raw[352+6:], 0xc041)
	if records, err := DecodeOverlays(raw, context); err == nil || records != nil {
		t.Fatal("invalid packed wire count accepted")
	}
	if encoded, err := EncodeOverlays([]lumps.Overlay{{FaceCountAndRenderOrder: 0xc041}}, context); err == nil || encoded != nil {
		t.Fatal("invalid packed record count accepted")
	}
}

func TestSpatialUnknownPrimitiveTypesRemainUnsupported(t *testing.T) {
	for _, context := range spatialContexts() {
		for _, kind := range []byte{2, 127, 255} {
			raw := spatialPrimitiveFixture()
			raw[10] = kind
			before := append([]byte(nil), raw...)
			if records, err := DecodePrimitives(raw, context); !errors.Is(err, ErrUnsupportedPrimitiveType) || records != nil {
				t.Fatalf("type %d decode: %+v %v", kind, records, err)
			}
			if !bytes.Equal(raw, before) {
				t.Fatal("unsupported decode mutated raw")
			}
			records := []lumps.Primitive{{Type: 0}, {Type: kind, Padding: 0xa5}}
			original := append([]lumps.Primitive(nil), records...)
			if encoded, err := EncodePrimitives(records, context); !errors.Is(err, ErrUnsupportedPrimitiveType) || encoded != nil {
				t.Fatalf("type %d encode: %v", kind, err)
			}
			if !reflect.DeepEqual(records, original) {
				t.Fatal("unsupported encode mutated records")
			}
		}
	}
}

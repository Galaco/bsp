package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"testing"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
	"github.com/go-gl/mathgl/mgl32"
)

func occlusionTestContext(version int32) format.Context {
	return format.Context{Profile: "sdk2013", BSPVersion: 20, LumpVersion: version, ByteOrder: binary.LittleEndian}
}

// Literal wire words, not production structs or encoders, are the fixture oracle.
func occlusionFixture(version int32) []byte {
	width := 36
	if version == 2 {
		width = 40
	}
	raw := make([]byte, 12+2*width+3*12+6*4)
	binary.LittleEndian.PutUint32(raw[0:4], 2)
	words := [2][10]uint32{
		{0xffffffef, 1, 2, 0xbfa00000, 0x40200000, 0xc0700000, 0x40880000, 0x40b00000, 0x40d80000, 13},
		{0x12345678, 0, 1, 0xc0f00000, 0x41040000, 0xc1100000, 0x41280000, 0x41340000, 0x41400000, 0xffffffe3},
	}
	for record, fields := range words {
		for field, word := range fields[:width/4] {
			offset := 4 + record*width + field*4
			binary.LittleEndian.PutUint32(raw[offset:offset+4], word)
		}
	}
	offset := 4 + 2*width
	binary.LittleEndian.PutUint32(raw[offset:offset+4], 3)
	offset += 4
	for _, fields := range [3][3]uint32{{0, 2, 7}, {2, 3, 3}, {5, 1, 12}} {
		for _, word := range fields {
			binary.LittleEndian.PutUint32(raw[offset:offset+4], word)
			offset += 4
		}
	}
	binary.LittleEndian.PutUint32(raw[offset:offset+4], 6)
	offset += 4
	for _, index := range []uint32{9, 2, 17, 4, 11, 0} {
		binary.LittleEndian.PutUint32(raw[offset:offset+4], index)
		offset += 4
	}
	return raw
}

func occlusionExpected(version int32) lumps.Occlusion {
	value := lumps.Occlusion{
		Data: []lumps.OcclusionData{
			{Flags: -17, FirstPoly: 1, PolyCount: 2, Mins: mgl32.Vec3{-1.25, 2.5, -3.75}, Maxs: mgl32.Vec3{4.25, 5.5, 6.75}},
			{Flags: 0x12345678, FirstPoly: 0, PolyCount: 1, Mins: mgl32.Vec3{-7.5, 8.25, -9}, Maxs: mgl32.Vec3{10.5, 11.25, 12}},
		},
		Polygons: []lumps.OcclusionPolyData{
			{FirstVertexIndex: 0, VertexCount: 2, PlaneNum: 7},
			{FirstVertexIndex: 2, VertexCount: 3, PlaneNum: 3},
			{FirstVertexIndex: 5, VertexCount: 1, PlaneNum: 12},
		},
		VertexIndices: []int32{9, 2, 17, 4, 11, 0},
	}
	if version == 2 {
		value.Data[0].Area = 13
		value.Data[1].Area = -29
		value.Data[0].HasArea = true
		value.Data[1].HasArea = true
	}
	return value
}

func TestOcclusionIndependentWireVersions(t *testing.T) {
	for _, version := range []int32{1, 2} {
		t.Run(fmt.Sprint(version), func(t *testing.T) {
			raw := occlusionFixture(version)
			context := occlusionTestContext(version)
			value, err := DecodeOcclusion(raw, context)
			if err != nil {
				t.Fatal(err)
			}
			if expected := occlusionExpected(version); !reflect.DeepEqual(value, expected) {
				t.Fatalf("decoded %+v, want %+v", value, expected)
			}
			encoded, err := EncodeOcclusion(occlusionExpected(version), context)
			if err != nil || !bytes.Equal(encoded, raw) {
				t.Fatalf("encode differs from independent wire fixture: %v", err)
			}
			if err := ValidateOcclusionReferences(value, 13, 18); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestOcclusionEveryTruncationAndTrailingByte(t *testing.T) {
	for _, version := range []int32{1, 2} {
		raw := occlusionFixture(version)
		for length := 0; length < len(raw); length++ {
			if _, err := DecodeOcclusion(raw[:length], occlusionTestContext(version)); err == nil {
				t.Fatalf("revision %d accepted prefix length %d", version, length)
			} else if !errors.Is(err, io.ErrUnexpectedEOF) {
				t.Fatalf("revision %d prefix %d error does not identify truncation: %v", version, length, err)
			}
		}
		if _, err := DecodeOcclusion(append(raw, 0x7f), occlusionTestContext(version)); err == nil {
			t.Fatalf("revision %d accepted trailing byte", version)
		}
	}
}

func TestOcclusionMaliciousCounts(t *testing.T) {
	for _, version := range []int32{1, 2} {
		width := 36
		if version == 2 {
			width = 40
		}
		for _, offset := range []int{0, 4 + 2*width, 8 + 2*width + 3*12} {
			for _, count := range []uint32{0xffffffff, 0x80000000, 0x7fffffff, 0x40000000, 1000} {
				raw := occlusionFixture(version)
				binary.LittleEndian.PutUint32(raw[offset:offset+4], count)
				if _, err := DecodeOcclusion(raw, occlusionTestContext(version)); err == nil {
					t.Fatalf("revision %d accepted count %#x at offset %d", version, count, offset)
				}
			}
		}
	}
	for _, version := range []int32{1, 2} {
		// No object bytes or later count words accompany the claimed object.
		if _, err := DecodeOcclusion([]byte{1, 0, 0, 0}, occlusionTestContext(version)); err == nil {
			t.Fatal("accepted missing array and later counts")
		}
	}
}

func TestOcclusionSpanValidation(t *testing.T) {
	cases := []struct {
		name         string
		field        int
		first, count int32
	}{
		{"negative-first-poly", 0, -1, 1},
		{"negative-poly-count", 0, 0, -1},
		{"past-polygons", 0, 3, 1},
		{"poly-addition-overflow", 0, math.MaxInt32, math.MaxInt32},
		{"poly-empty-past-end", 0, 4, 0},
		{"negative-first-index", 1, -1, 1},
		{"negative-index-count", 1, 0, -1},
		{"past-indices", 1, 6, 1},
		{"index-addition-overflow", 1, math.MaxInt32, math.MaxInt32},
		{"index-empty-past-end", 1, 7, 0},
	}
	for _, version := range []int32{1, 2} {
		for _, test := range cases {
			t.Run(fmt.Sprintf("%d/%s", version, test.name), func(t *testing.T) {
				value := occlusionExpected(version)
				raw := occlusionFixture(version)
				offset := 8
				if test.field == 0 {
					value.Data[0].FirstPoly, value.Data[0].PolyCount = test.first, test.count
				} else {
					value.Polygons[0].FirstVertexIndex, value.Polygons[0].VertexCount = test.first, test.count
					width := 36
					if version == 2 {
						width = 40
					}
					offset = 8 + 2*width
				}
				binary.LittleEndian.PutUint32(raw[offset:offset+4], uint32(test.first))
				binary.LittleEndian.PutUint32(raw[offset+4:offset+8], uint32(test.count))
				if _, err := DecodeOcclusion(raw, occlusionTestContext(version)); err == nil {
					t.Fatal("decode accepted invalid span")
				}
				if _, err := EncodeOcclusion(value, occlusionTestContext(version)); err == nil {
					t.Fatal("encode accepted invalid span")
				}
				if err := ValidateOcclusionReferences(value, 13, 18); err == nil {
					t.Fatal("reference validator accepted invalid span")
				}
			})
		}
		value := occlusionExpected(version)
		value.Data[0].FirstPoly, value.Data[0].PolyCount = 3, 0
		value.Polygons[0].FirstVertexIndex, value.Polygons[0].VertexCount = 6, 0
		if _, err := EncodeOcclusion(value, occlusionTestContext(version)); err != nil {
			t.Fatalf("empty spans ending at exact bounds must be valid: %v", err)
		}
	}
}

func TestOcclusionReferenceBoundaries(t *testing.T) {
	for _, test := range []struct {
		name             string
		planes, vertices int
		plane, vertex    int32
	}{
		{"negative-plane-count", -1, 18, 7, 9},
		{"negative-vertex-count", 13, -1, 7, 9},
		{"plane-exclusive-bound", 13, 18, 13, 9},
		{"vertex-exclusive-bound", 13, 18, 7, 18},
		{"negative-plane", 13, 18, -1, 9},
		{"negative-vertex", 13, 18, 7, -1},
		{"max-plane", 13, 18, math.MaxInt32, 9},
		{"max-vertex", 13, 18, 7, math.MaxInt32},
		{"no-planes", 0, 18, 0, 9},
		{"no-vertices", 13, 0, 7, 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			value := occlusionExpected(2)
			// Leave every polygon and index unselected: references still matter.
			value.Data = nil
			value.Polygons[0].PlaneNum = test.plane
			value.VertexIndices[0] = test.vertex
			if err := ValidateOcclusionReferences(value, test.planes, test.vertices); err == nil {
				t.Fatal("accepted invalid external reference")
			}
		})
	}
	if err := ValidateOcclusionReferences(lumps.Occlusion{}, 0, 0); err != nil {
		t.Fatal(err)
	}
}

func TestOcclusionRevisionOneRejectsAreaEdits(t *testing.T) {
	for _, edit := range []struct {
		area    int32
		hasArea bool
	}{{0, true}, {19, false}, {-1, false}, {13, true}} {
		value := occlusionExpected(1)
		value.Data[1].Area, value.Data[1].HasArea = edit.area, edit.hasArea
		if _, err := EncodeOcclusion(value, occlusionTestContext(1)); err == nil {
			t.Fatalf("accepted unrepresentable area edit %+v", edit)
		}
	}
	value := occlusionExpected(2)
	value.Data[0].HasArea = false
	raw, err := EncodeOcclusion(value, occlusionTestContext(2))
	if err != nil || !bytes.Equal(raw, occlusionFixture(2)) {
		t.Fatalf("revision 2 must encode area regardless of presence marker: %v", err)
	}
}

func TestOcclusionContexts(t *testing.T) {
	valid := []format.Context{
		{Profile: "sdk2013", BSPVersion: 19, LumpVersion: 1},
		{Profile: "sdk2013", BSPVersion: 20, LumpVersion: 2},
		{Profile: "csgo", BSPVersion: 21, LumpVersion: 1},
		{Profile: "csgo", BSPVersion: 21, LumpVersion: 2},
		{Profile: "l4d2", BSPVersion: 21, LumpVersion: 1},
		{Profile: "l4d2", BSPVersion: 21, LumpVersion: 2},
		{Profile: "reactivedrop", BSPVersion: 21, LumpVersion: 2},
	}
	for _, context := range valid {
		if _, err := DecodeOcclusion(occlusionFixture(context.LumpVersion), context); err != nil {
			t.Fatalf("valid context %+v: %v", context, err)
		}
	}
	invalid := []format.Context{
		{},
		{Profile: "sdk2013", BSPVersion: 20, LumpVersion: 0},
		{Profile: "sdk2013", BSPVersion: 20, LumpVersion: -1},
		{Profile: "sdk2013", BSPVersion: 20, LumpVersion: 3},
		{Profile: "sdk2013", BSPVersion: 21, LumpVersion: 2},
		{Profile: "csgo", BSPVersion: 20, LumpVersion: 2},
		{Profile: "l4d2", BSPVersion: 19, LumpVersion: 2},
		{Profile: "guessed", BSPVersion: 20, LumpVersion: 2},
		{Profile: "sdk2013", BSPVersion: 20, LumpVersion: 2, ByteOrder: binary.BigEndian},
		{Profile: "reactivedrop", BSPVersion: 21, LumpVersion: 0},
		{Profile: "reactivedrop", BSPVersion: 21, LumpVersion: 1},
		{Profile: "reactivedrop", BSPVersion: 21, LumpVersion: 3},
		{Profile: "reactivedrop", BSPVersion: 19, LumpVersion: 2},
		{Profile: "reactivedrop", BSPVersion: 20, LumpVersion: 2},
		{Profile: "reactivedrop", BSPVersion: 21, LumpVersion: 2, ByteOrder: binary.BigEndian},
		{Profile: "reactivedrop", BSPVersion: 21, LumpVersion: 2, HeaderLayout: format.L4D2Header},
	}
	for _, context := range invalid {
		if _, err := DecodeOcclusion(occlusionFixture(2), context); err == nil {
			t.Fatalf("decode accepted context %+v", context)
		}
		if _, err := EncodeOcclusion(occlusionExpected(2), context); err == nil {
			t.Fatalf("encode accepted context %+v", context)
		}
	}
	if _, err := DecodeOcclusion(occlusionFixture(1), occlusionTestContext(2)); err == nil {
		t.Fatal("decoded v1 bytes as v2")
	}
	if _, err := DecodeOcclusion(occlusionFixture(2), occlusionTestContext(1)); err == nil {
		t.Fatal("decoded v2 bytes as v1")
	}
}

func TestOcclusionEmptyAndDerivedCounts(t *testing.T) {
	for _, version := range []int32{1, 2} {
		raw := make([]byte, 12)
		value, err := DecodeOcclusion(raw, occlusionTestContext(version))
		if err != nil || len(value.Data)+len(value.Polygons)+len(value.VertexIndices) != 0 {
			t.Fatalf("empty counted payload: %v", err)
		}
		encoded, err := EncodeOcclusion(lumps.Occlusion{}, occlusionTestContext(version))
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("empty payload encode: %v", err)
		}
		// Preserve records even when no occluder references them; count is now 1.
		value = occlusionExpected(version)
		value.Data = value.Data[:1]
		encoded, err = EncodeOcclusion(value, occlusionTestContext(version))
		if err != nil {
			t.Fatal(err)
		}
		if binary.LittleEndian.Uint32(encoded[:4]) != 1 {
			t.Fatal("object count was not derived from edited slice")
		}
		decoded, err := DecodeOcclusion(encoded, occlusionTestContext(version))
		if err != nil || !reflect.DeepEqual(decoded, value) {
			t.Fatalf("unselected records were not preserved: %v", err)
		}
	}
}

func TestOcclusionFloatBitsPreserved(t *testing.T) {
	for _, version := range []int32{1, 2} {
		raw := occlusionFixture(version)
		binary.LittleEndian.PutUint32(raw[16:20], 0x80000000) // negative zero
		binary.LittleEndian.PutUint32(raw[20:24], 0x7fc12345) // NaN payload
		value, err := DecodeOcclusion(raw, occlusionTestContext(version))
		if err != nil {
			t.Fatal(err)
		}
		if math.Float32bits(value.Data[0].Mins[0]) != 0x80000000 || math.Float32bits(value.Data[0].Mins[1]) != 0x7fc12345 {
			t.Fatal("float bit patterns changed during decoding")
		}
		encoded, err := EncodeOcclusion(value, occlusionTestContext(version))
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("float bit patterns changed during encoding: %v", err)
		}
	}
}

func TestOcclusionEncodedSizeBoundaries(t *testing.T) {
	const maxSize = int64(1<<31 - 1)
	for _, width := range []int{36, 40} {
		maxRecords := int((maxSize - 12) / int64(width))
		if size, err := occlusionEncodedSize(maxRecords, 0, 0, width); err != nil || int64(size) != 12+int64(maxRecords)*int64(width) {
			t.Fatalf("valid maximum object count size=%d err=%v", size, err)
		}
		if _, err := occlusionEncodedSize(maxRecords+1, 0, 0, width); err == nil {
			t.Fatal("accepted object byte-size overflow")
		}
		if _, err := occlusionEncodedSize(maxRecords, 4, 0, width); err == nil {
			t.Fatal("accepted combined array byte-size overflow")
		}
	}
	for _, counts := range [][3]int{{-1, 0, 0}, {0, -1, 0}, {0, 0, -1}, {math.MaxInt32, 0, 0}, {0, math.MaxInt32, 0}, {0, 0, math.MaxInt32}} {
		if _, err := occlusionEncodedSize(counts[0], counts[1], counts[2], 40); err == nil {
			t.Fatalf("accepted invalid counts %v", counts)
		}
	}
	overflowCount := maxSize + 1
	if uint64(^uint(0)>>1) >= uint64(overflowCount) {
		for _, counts := range [][3]int{{int(overflowCount), 0, 0}, {0, int(overflowCount), 0}, {0, 0, int(overflowCount)}} {
			if _, err := occlusionEncodedSize(counts[0], counts[1], counts[2], 40); err == nil {
				t.Fatalf("accepted signed count overflow %v", counts)
			}
		}
	}
}

func TestReactiveDropOcclusionVersionTwoWire(t *testing.T) {
	context := format.Context{Profile: "reactivedrop", BSPVersion: 21, LumpVersion: 2}
	raw, want := occlusionFixture(2), occlusionExpected(2)
	got, err := DecodeOcclusion(raw, context)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("RD v2 occlusion: got %+v err=%v want %+v", got, err, want)
	}
	encoded, err := EncodeOcclusion(want, context)
	if err != nil || !bytes.Equal(encoded, raw) {
		t.Fatalf("RD v2 occlusion independent encoding: %v", err)
	}
	got.Data[1].Area = -1234
	mutated := append([]byte(nil), raw...)
	binary.LittleEndian.PutUint32(mutated[80:], 0xfffffb2e)
	encoded, err = EncodeOcclusion(got, context)
	if err != nil || !bytes.Equal(encoded, mutated) {
		t.Fatalf("RD occlusion area edit changed other fields: %v", err)
	}
	for n := 0; n < len(raw); n++ {
		if _, err := DecodeOcclusion(raw[:n], context); err == nil {
			t.Fatalf("RD occlusion accepted partial length %d", n)
		}
	}
	if _, err := DecodeOcclusion(append(raw, 0xd3), context); err == nil {
		t.Fatal("RD occlusion accepted trailing byte")
	}
	if _, err := DecodeOcclusion(occlusionFixture(1), context); err == nil {
		t.Fatal("RD occlusion treated 36-byte v1 records as v2")
	}
}

func FuzzOcclusionCountedSections(f *testing.F) {
	for _, version := range []int32{1, 2} {
		raw := occlusionFixture(version)
		f.Add(raw, version == 2)
		f.Add(make([]byte, 12), version == 2)
		f.Add(raw[:len(raw)-1], version == 2)
		f.Add(append(append([]byte(nil), raw...), 0xa5), version == 2)
		width := 36
		if version == 2 {
			width = 40
		}
		for _, offset := range []int{0, 4 + 2*width, 8 + 2*width + 3*12} {
			for _, count := range []uint32{0xffffffff, 0x7fffffff} {
				malformed := append([]byte(nil), raw...)
				binary.LittleEndian.PutUint32(malformed[offset:], count)
				f.Add(malformed, version == 2)
			}
		}
	}
	f.Add([]byte{}, false)
	f.Fuzz(func(t *testing.T, raw []byte, revisionTwo bool) {
		if len(raw) > 4096 {
			return
		}
		raw = append([]byte(nil), raw...)
		version, width := int32(1), 36
		if revisionTwo {
			version, width = 2, 40
		}
		context := occlusionTestContext(version)
		initial := occlusionFixture(version)
		var adapter Occlusion
		adapter.SetContext(context)
		if err := adapter.Decode(initial); err != nil {
			t.Fatal(err)
		}
		before := append([]byte(nil), raw...)
		err := adapter.Decode(raw)
		if !bytes.Equal(raw, before) {
			t.Fatal("decode mutated source")
		}
		encoded, encodeErr := adapter.Encode()
		if err != nil {
			if encodeErr != nil || adapter.Length() != len(initial) || !bytes.Equal(encoded, initial) {
				t.Fatalf("failed counted decode changed committed state: %v", encodeErr)
			}
			return
		}
		value := adapter.GetData()
		if len(raw) != 12+width*len(value.Data)+12*len(value.Polygons)+4*len(value.VertexIndices) {
			t.Fatal("counted sections did not consume exactly the selected revision")
		}
		for i, record := range value.Data {
			if record.FirstPoly < 0 || record.PolyCount < 0 || int64(record.FirstPoly)+int64(record.PolyCount) > int64(len(value.Polygons)) {
				t.Fatalf("occluder %d has an invalid polygon span", i)
			}
		}
		for i, record := range value.Polygons {
			if record.FirstVertexIndex < 0 || record.VertexCount < 0 || int64(record.FirstVertexIndex)+int64(record.VertexCount) > int64(len(value.VertexIndices)) {
				t.Fatalf("polygon %d has an invalid vertex-index span", i)
			}
		}
		if encodeErr != nil || len(encoded) != len(raw) || !bytes.Equal(encoded, before) || !bytes.Equal(raw, before) {
			t.Fatalf("occlusion wire/float bits changed on export: %v", encodeErr)
		}
		raw[len(raw)-1] ^= 0xff
		encoded[0] ^= 0xff
		repeated, err := adapter.Encode()
		if err != nil || !bytes.Equal(repeated, before) {
			t.Fatalf("occlusion snapshot borrowed source or export bytes: %v", err)
		}
	})
}

package codec

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
	"github.com/galaco/bsp/primitives/common"
)

func leafTestContext(version int32) format.Context {
	return format.Context{Profile: "sdk2013", BSPVersion: 20, LumpVersion: version, LeafVersion: version}
}

// These records are authored independently of the encoder. Distinct fields,
// signed exponents and nonzero padding catch wrong widths and silent discard.
func leafWireRecords(version int32) []byte {
	first := []byte{
		1, 0, 0, 0, 2, 0, 61, 92,
		3, 0, 7, 0, 14, 0, 13, 0, 36, 0, 51, 0,
		32, 0, 11, 0, 21, 0, 76, 0, 109, 2,
	}
	second := []byte{
		254, 255, 255, 255, 255, 255, 3, 130,
		255, 255, 254, 255, 253, 255, 4, 0, 5, 0, 6, 0,
		52, 18, 120, 86, 188, 154, 240, 222, 255, 255,
	}
	if version == 0 {
		first = append(first,
			1, 2, 3, 255, 4, 5, 6, 254, 7, 8, 9, 253,
			10, 11, 12, 252, 13, 14, 15, 251, 16, 17, 18, 250)
		second = append(second,
			21, 22, 23, 1, 24, 25, 26, 2, 27, 28, 29, 3,
			30, 31, 32, 4, 33, 34, 35, 5, 36, 37, 38, 6)
	}
	first = append(first, 0xab, 0xcd)
	second = append(second, 0x98, 0x76)
	return append(first, second...)
}

func leafExpectedRecords(version int32) []lumps.Leaf {
	data := []lumps.Leaf{
		{
			Contents: 1, Cluster: 2, BitField: 0x5c3d,
			Mins: [3]int16{3, 7, 14}, Maxs: [3]int16{13, 36, 51},
			FirstLeafFace: 32, NumLeafFaces: 11, FirstLeafBrush: 21,
			NumLeafBrushes: 76, LeafWaterDataID: 621, Padding: [2]byte{0xab, 0xcd},
		},
		{
			Contents: -2, Cluster: -1, BitField: 0x8203,
			Mins: [3]int16{-1, -2, -3}, Maxs: [3]int16{4, 5, 6},
			FirstLeafFace: 0x1234, NumLeafFaces: 0x5678, FirstLeafBrush: 0x9abc,
			NumLeafBrushes: 0xdef0, LeafWaterDataID: -1, Padding: [2]byte{0x98, 0x76},
		},
	}
	if version == 0 {
		data[0].HasInlineLighting, data[1].HasInlineLighting = true, true
		data[0].LightSample.Color = [6]common.ColorRGBExponent32{
			{R: 1, G: 2, B: 3, Exponent: -1}, {R: 4, G: 5, B: 6, Exponent: -2},
			{R: 7, G: 8, B: 9, Exponent: -3}, {R: 10, G: 11, B: 12, Exponent: -4},
			{R: 13, G: 14, B: 15, Exponent: -5}, {R: 16, G: 17, B: 18, Exponent: -6},
		}
		data[1].LightSample.Color = [6]common.ColorRGBExponent32{
			{R: 21, G: 22, B: 23, Exponent: 1}, {R: 24, G: 25, B: 26, Exponent: 2},
			{R: 27, G: 28, B: 29, Exponent: 3}, {R: 30, G: 31, B: 32, Exponent: 4},
			{R: 33, G: 34, B: 35, Exponent: 5}, {R: 36, G: 37, B: 38, Exponent: 6},
		}
	}
	return data
}

func TestLeavesEveryFieldAndPadding(t *testing.T) {
	for _, version := range []int32{0, 1} {
		context := leafTestContext(version) // Same BSP/profile, differing descriptors.
		raw, want := leafWireRecords(version), leafExpectedRecords(version)
		got, err := DecodeLeaves(raw, context)
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("version %d: got %+v, err %v; want %+v", version, got, err, want)
		}
		encoded, err := EncodeLeaves(want, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("version %d authored records: %x, %v; want %x", version, encoded, err, raw)
		}
		encoded, err = EncodeLeaves(got, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("version %d round trip: %x, %v", version, encoded, err)
		}
		// A supported mutation must preserve all unedited bytes including padding.
		got[1].Cluster = -32768
		encoded, err = EncodeLeaves(got, context)
		stride, _ := LeafRecordSize(context)
		mutated := append([]byte(nil), raw...)
		mutated[stride+4], mutated[stride+5] = 0, 0x80
		if err != nil || !bytes.Equal(encoded, mutated) {
			t.Fatalf("version %d mutation: %x, %v", version, encoded, err)
		}
		context.LeafVersion = 99 // Sibling ambient context never selects leaf width.
		if size, err := LeafRecordSize(context); err != nil || size != stride {
			t.Fatalf("leaf width follows sibling context: %d, %v", size, err)
		}
	}
}

func TestLeavesRejectEveryPartialRecord(t *testing.T) {
	for _, version := range []int32{0, 1} {
		context := leafTestContext(version)
		raw := leafWireRecords(version)
		stride, _ := LeafRecordSize(context)
		for cut := 1; cut < len(raw); cut++ {
			if cut%stride == 0 {
				continue
			}
			if got, err := DecodeLeaves(raw[:cut], context); err == nil || got != nil {
				t.Fatalf("version %d accepted partial record length %d", version, cut)
			}
		}
		for trailing := 1; trailing < stride; trailing++ {
			if _, err := DecodeLeaves(append(append([]byte(nil), raw...), make([]byte, trailing)...), context); err == nil {
				t.Fatalf("version %d accepted %d trailing bytes", version, trailing)
			}
		}
	}
}

func TestLeavesCountBoundaries(t *testing.T) {
	for _, version := range []int32{0, 1} {
		context := leafTestContext(version)
		stride, _ := LeafRecordSize(context)
		for _, count := range []int{0, lumps.MaxMapLeafs} {
			raw := make([]byte, count*stride)
			data, err := DecodeLeaves(raw, context)
			if err != nil || len(data) != count {
				t.Fatalf("count %d: %v", count, err)
			}
			encoded, err := EncodeLeaves(data, context)
			if err != nil || !bytes.Equal(encoded, raw) {
				t.Fatalf("count %d encode: %v", count, err)
			}
		}
		if _, err := DecodeLeaves(make([]byte, (lumps.MaxMapLeafs+1)*stride), context); err == nil {
			t.Fatal("decoded too many leaves")
		}
		if _, err := EncodeLeaves(make([]lumps.Leaf, lumps.MaxMapLeafs+1), context); err == nil {
			t.Fatal("encoded too many leaves")
		}
	}
}

func TestLeavesRejectNonrepresentableFields(t *testing.T) {
	for _, version := range []int32{0, 1} {
		data := leafExpectedRecords(version)
		data[1].HasInlineLighting = !data[1].HasInlineLighting
		if raw, err := EncodeLeaves(data, leafTestContext(version)); err == nil || raw != nil {
			t.Fatalf("version %d accepted mismatched inline-lighting flag", version)
		}
	}
	data := leafExpectedRecords(1)
	data[1].LightSample.Color[5].Exponent = -128
	if raw, err := EncodeLeaves(data, leafTestContext(1)); err == nil || raw != nil {
		t.Fatal("version 1 discarded inline-lighting field")
	}
}

func TestLeavesRejectWrongContexts(t *testing.T) {
	contexts := []format.Context{
		{}, {Profile: "future", BSPVersion: 20},
		{Profile: "sdk2013", BSPVersion: 21}, {Profile: "csgo", BSPVersion: 20},
		{Profile: "sdk2013", BSPVersion: 20, ByteOrder: binary.BigEndian},
	}
	for _, version := range []int32{-1, 2, 19, 20, 21} {
		context := leafTestContext(version)
		contexts = append(contexts, context)
	}
	for _, context := range contexts {
		if _, err := LeafRecordSize(context); err == nil {
			t.Fatalf("accepted context %+v", context)
		}
		if _, err := DecodeLeaves(nil, context); err == nil {
			t.Fatalf("decoded context %+v", context)
		}
		if _, err := EncodeLeaves(nil, context); err == nil {
			t.Fatalf("encoded context %+v", context)
		}
	}
	for _, profile := range []string{"csgo", "l4d2"} {
		for _, version := range []int32{0, 1} {
			context := leafTestContext(version)
			context.Profile, context.BSPVersion = profile, 21
			if _, err := DecodeLeaves(leafWireRecords(version), context); err != nil {
				t.Fatal(err)
			}
		}
	}
}

// RD's pinned LoadLeafs explicitly accepts these two descriptor revisions;
// BSP21 alone does not remove the old inline cube or its alignment bytes.
func TestReactiveDropLeavesWireAndContext(t *testing.T) {
	for _, version := range []int32{0, 1} {
		context := format.Context{Profile: "reactivedrop", BSPVersion: 21, LumpVersion: version, ByteOrder: binary.LittleEndian}
		raw, want := leafWireRecords(version), leafExpectedRecords(version)
		got, err := DecodeLeaves(raw, context)
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("RD leaf-v%d: got %+v err=%v want %+v", version, got, err, want)
		}
		encoded, err := EncodeLeaves(want, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("RD leaf-v%d authored encoding: %v", version, err)
		}
		stride := 32
		if version == 0 {
			stride = 56
		}
		got[1].Padding = [2]byte{0x12, 0xef}
		mutated := append([]byte(nil), raw...)
		mutated[2*stride-2], mutated[2*stride-1] = 0x12, 0xef
		encoded, err = EncodeLeaves(got, context)
		if err != nil || !bytes.Equal(encoded, mutated) {
			t.Fatalf("RD leaf-v%d padding edit changed other fields: %v", version, err)
		}
		for n := 1; n < len(raw); n++ {
			if n == stride {
				continue
			}
			if value, err := DecodeLeaves(raw[:n], context); err == nil || value != nil {
				t.Fatalf("RD leaf-v%d accepted partial length %d", version, n)
			}
		}
	}
	for _, context := range []format.Context{
		{Profile: "reactivedrop", BSPVersion: 19, LumpVersion: 1},
		{Profile: "reactivedrop", BSPVersion: 20, LumpVersion: 1},
		{Profile: "reactivedrop", BSPVersion: 22, LumpVersion: 1},
		{Profile: "reactivedrop", BSPVersion: 21, LumpVersion: -1},
		{Profile: "reactivedrop", BSPVersion: 21, LumpVersion: 2},
		{Profile: "reactivedrop", BSPVersion: 21, LumpVersion: 1, ByteOrder: binary.BigEndian},
		{Profile: "reactivedrop", BSPVersion: 21, LumpVersion: 1, HeaderLayout: format.L4D2Header},
	} {
		if value, err := DecodeLeaves(leafWireRecords(1), context); err == nil || value != nil {
			t.Fatalf("RD leaves accepted context %+v", context)
		}
		if raw, err := EncodeLeaves(leafExpectedRecords(1), context); err == nil || raw != nil {
			t.Fatalf("RD leaves encoded context %+v", context)
		}
	}
}

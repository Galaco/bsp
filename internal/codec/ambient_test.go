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

func ambientTestContext(legacy bool) format.Context {
	return format.Context{Profile: "sdk2013", BSPVersion: 20, LumpVersion: 1, LeafVersion: 1, LegacyAmbient: legacy}
}

func ambientWireRecords(positioned bool) []byte {
	first := []byte{
		1, 2, 3, 255, 4, 5, 6, 254, 7, 8, 9, 253,
		10, 11, 12, 252, 13, 14, 15, 251, 16, 17, 18, 250,
	}
	second := []byte{
		21, 22, 23, 1, 24, 25, 26, 2, 27, 28, 29, 3,
		30, 31, 32, 4, 33, 34, 35, 5, 36, 37, 38, 128,
	}
	if positioned {
		first = append(first, 0x12, 0x34, 0x56, 0xab)
		second = append(second, 0x98, 0x76, 0x54, 0xcd)
	}
	return append(first, second...)
}

func ambientExpectedRecords(positioned bool) []lumps.AmbientSample {
	data := []lumps.AmbientSample{
		{Cube: common.CompressedLightCube{Color: [6]common.ColorRGBExponent32{
			{R: 1, G: 2, B: 3, Exponent: -1}, {R: 4, G: 5, B: 6, Exponent: -2},
			{R: 7, G: 8, B: 9, Exponent: -3}, {R: 10, G: 11, B: 12, Exponent: -4},
			{R: 13, G: 14, B: 15, Exponent: -5}, {R: 16, G: 17, B: 18, Exponent: -6},
		}}},
		{Cube: common.CompressedLightCube{Color: [6]common.ColorRGBExponent32{
			{R: 21, G: 22, B: 23, Exponent: 1}, {R: 24, G: 25, B: 26, Exponent: 2},
			{R: 27, G: 28, B: 29, Exponent: 3}, {R: 30, G: 31, B: 32, Exponent: 4},
			{R: 33, G: 34, B: 35, Exponent: 5}, {R: 36, G: 37, B: 38, Exponent: -128},
		}}},
	}
	if positioned {
		data[0].X, data[0].Y, data[0].Z, data[0].Pad, data[0].HasPosition = 0x12, 0x34, 0x56, 0xab, true
		data[1].X, data[1].Y, data[1].Z, data[1].Pad, data[1].HasPosition = 0x98, 0x76, 0x54, 0xcd, true
	}
	return data
}

func TestAmbientSamplesEveryFieldAndPadding(t *testing.T) {
	for _, legacy := range []bool{false, true} {
		context := ambientTestContext(legacy)
		versions := []int32{0, 1}
		if legacy {
			versions = append(versions, -123, 99999)
		}
		for _, version := range versions {
			context.LumpVersion = version
			raw, want := ambientWireRecords(!legacy), ambientExpectedRecords(!legacy)
			got, err := DecodeAmbientSamples(raw, context)
			if err != nil || !reflect.DeepEqual(got, want) {
				t.Fatalf("legacy %v version %d: %+v, %v", legacy, version, got, err)
			}
			encoded, err := EncodeAmbientSamples(want, context)
			if err != nil || !bytes.Equal(encoded, raw) {
				t.Fatalf("authored ambient records: %x, %v", encoded, err)
			}
			encoded, err = EncodeAmbientSamples(got, context)
			if err != nil || !bytes.Equal(encoded, raw) {
				t.Fatalf("ambient round trip: %x, %v", encoded, err)
			}
			got[1].Cube.Color[5].Exponent = 127
			encoded, err = EncodeAmbientSamples(got, context)
			mutated := append([]byte(nil), raw...)
			stride := 28
			if legacy {
				stride = 24
			}
			mutated[stride+23] = 127
			if err != nil || !bytes.Equal(encoded, mutated) {
				t.Fatalf("ambient mutation: %x, %v", encoded, err)
			}
		}
	}
}

func TestAmbientIndicesEveryField(t *testing.T) {
	context := ambientTestContext(false)
	context.LumpVersion = 0
	raw := []byte{0x34, 0x12, 0x78, 0x56, 0xff, 0xff, 0xbc, 0x9a}
	want := []lumps.AmbientIndex{{AmbientSampleCount: 0x1234, FirstAmbientSample: 0x5678}, {AmbientSampleCount: 65535, FirstAmbientSample: 0x9abc}}
	got, err := DecodeAmbientIndices(raw, context)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("indices: %+v, %v", got, err)
	}
	encoded, err := EncodeAmbientIndices(want, context)
	if err != nil || !bytes.Equal(encoded, raw) {
		t.Fatalf("authored indices: %x, %v", encoded, err)
	}
	got[1].FirstAmbientSample = 65535
	encoded, err = EncodeAmbientIndices(got, context)
	if err != nil || !bytes.Equal(encoded, []byte{0x34, 0x12, 0x78, 0x56, 0xff, 0xff, 0xff, 0xff}) {
		t.Fatalf("index mutation: %x, %v", encoded, err)
	}
}

func TestAmbientRejectsEveryPartialRecord(t *testing.T) {
	for _, legacy := range []bool{false, true} {
		context := ambientTestContext(legacy)
		raw := ambientWireRecords(!legacy)
		stride := 28
		if legacy {
			stride = 24
		}
		for cut := 1; cut < len(raw); cut++ {
			if cut%stride == 0 {
				continue
			}
			if got, err := DecodeAmbientSamples(raw[:cut], context); err == nil || got != nil {
				t.Fatalf("accepted ambient length %d", cut)
			}
		}
		for trailing := 1; trailing < stride; trailing++ {
			if _, err := DecodeAmbientSamples(append(append([]byte(nil), raw...), make([]byte, trailing)...), context); err == nil {
				t.Fatalf("accepted %d trailing sample bytes", trailing)
			}
		}
	}
	context := ambientTestContext(false)
	context.LumpVersion = 0
	for length := 1; length < 12; length++ {
		if length%4 == 0 {
			continue
		}
		if got, err := DecodeAmbientIndices(make([]byte, length), context); err == nil || got != nil {
			t.Fatalf("accepted index length %d", length)
		}
	}
}

func TestAmbientUsesSharedContextNotDivisibilityOrOwnRevision(t *testing.T) {
	// 168 bytes are divisible by BOTH widths. Neither that ambiguity nor a
	// sample revision (even revision 1 in legacy mode) overrides sibling context.
	raw := make([]byte, 168)
	raw[24], raw[27], raw[28] = 42, 99, 123
	for _, legacy := range []bool{false, true} {
		context := ambientTestContext(legacy)
		got, err := DecodeAmbientSamples(raw, context)
		count := 6
		if legacy {
			count = 7
		}
		if err != nil || len(got) != count || got[0].HasPosition == legacy {
			t.Fatalf("layout selected by bytes: legacy=%v got=%+v err=%v", legacy, got, err)
		}
		if legacy {
			if got[1].Cube.Color[0].R != 42 || got[0].X != 0 {
				t.Fatal("legacy interpreted positioned bytes")
			}
		} else if got[0].X != 42 || got[0].Pad != 99 || got[1].Cube.Color[0].R != 123 {
			t.Fatal("indexed interpreted legacy bytes")
		}
	}
	// A present HDR index with absent LDR index does not change LegacyAmbient.
	context := ambientTestContext(true)
	context.LumpVersion = 0
	if _, err := DecodeAmbientIndices([]byte{1, 0, 0, 0}, context); err == nil {
		t.Fatal("HDR index independently switched legacy layout")
	}
	if _, err := EncodeAmbientIndices([]lumps.AmbientIndex{{AmbientSampleCount: 1}}, context); err == nil {
		t.Fatal("encoded ignored legacy HDR indices")
	}
}

func TestAmbientLeafZeroAndAbsentChannels(t *testing.T) {
	for _, legacy := range []bool{false, true} {
		context := ambientTestContext(legacy)
		context.LeafVersion, context.LumpVersion = 0, 999
		if got, err := DecodeAmbientSamples(nil, context); err != nil || len(got) != 0 {
			t.Fatalf("empty leaf0 samples: %v", err)
		}
		if got, err := DecodeAmbientIndices(nil, context); err != nil || len(got) != 0 {
			t.Fatalf("empty leaf0 indices: %v", err)
		}
		if _, err := EncodeAmbientSamples(nil, context); err != nil {
			t.Fatal(err)
		}
		if _, err := EncodeAmbientIndices(nil, context); err != nil {
			t.Fatal(err)
		}
		if _, err := DecodeAmbientSamples(ambientWireRecords(false), context); err == nil {
			t.Fatal("leaf0 interpreted separate samples")
		}
		if _, err := DecodeAmbientIndices([]byte{1, 0, 0, 0}, context); err == nil {
			t.Fatal("leaf0 interpreted separate indices")
		}
		if _, err := EncodeAmbientSamples(ambientExpectedRecords(false), context); err == nil {
			t.Fatal("leaf0 encoded separate samples")
		}
		if _, err := EncodeAmbientIndices([]lumps.AmbientIndex{{}}, context); err == nil {
			t.Fatal("leaf0 encoded separate indices")
		}
		if err := ValidateAmbient(nil, nil, 2, context); err != nil {
			t.Fatal(err)
		}
	}
	for _, legacy := range []bool{false, true} {
		context := ambientTestContext(legacy)
		if err := ValidateAmbient(nil, nil, 2, context); err != nil {
			t.Fatalf("absent channel: %v", err)
		}
		if err := ValidateAmbient(nil, ambientExpectedRecords(!legacy), 2, context); (err == nil) != legacy {
			t.Fatalf("unindexed samples legacy=%v: %v", legacy, err)
		}
	}
}

func TestAmbientRejectsNonrepresentableFields(t *testing.T) {
	for _, field := range []string{"X", "Y", "Z", "Pad", "HasPosition"} {
		data := ambientExpectedRecords(false)
		switch field {
		case "X":
			data[1].X = 1
		case "Y":
			data[1].Y = 1
		case "Z":
			data[1].Z = 1
		case "Pad":
			data[1].Pad = 1
		case "HasPosition":
			data[1].HasPosition = true
		}
		if raw, err := EncodeAmbientSamples(data, ambientTestContext(true)); err == nil || raw != nil {
			t.Fatalf("discarded legacy %s", field)
		}
		if err := ValidateAmbient(nil, data, 2, ambientTestContext(true)); err == nil {
			t.Fatalf("validation discarded legacy %s", field)
		}
	}
	data := ambientExpectedRecords(true)
	data[1].HasPosition = false
	if _, err := EncodeAmbientSamples(data, ambientTestContext(false)); err == nil {
		t.Fatal("accepted positioned layout mismatch")
	}
}

func TestAmbientCrossLumpSpansAndLeafReferences(t *testing.T) {
	context := ambientTestContext(false)
	samples := ambientExpectedRecords(true)
	for _, tc := range []struct {
		name    string
		indices []lumps.AmbientIndex
		samples []lumps.AmbientSample
		leaves  int
		valid   bool
	}{
		{"two-spans", []lumps.AmbientIndex{{AmbientSampleCount: 1, FirstAmbientSample: 0}, {AmbientSampleCount: 1, FirstAmbientSample: 1}}, samples, 2, true},
		{"shared-spans", []lumps.AmbientIndex{{AmbientSampleCount: 2, FirstAmbientSample: 0}, {AmbientSampleCount: 1, FirstAmbientSample: 1}}, samples, 2, true},
		{"leaf-reference", []lumps.AmbientIndex{{AmbientSampleCount: 0, FirstAmbientSample: 1}, {AmbientSampleCount: 2, FirstAmbientSample: 0}}, samples, 2, true},
		{"leaf-reference-no-samples", []lumps.AmbientIndex{{AmbientSampleCount: 0, FirstAmbientSample: 1}, {AmbientSampleCount: 0, FirstAmbientSample: 0}}, nil, 2, true},
		{"reference-at-leaf-count", []lumps.AmbientIndex{{AmbientSampleCount: 0, FirstAmbientSample: 2}, {AmbientSampleCount: 2, FirstAmbientSample: 0}}, samples, 2, false},
		{"reference-max", []lumps.AmbientIndex{{AmbientSampleCount: 0, FirstAmbientSample: 65535}, {AmbientSampleCount: 2, FirstAmbientSample: 0}}, samples, 2, false},
		{"missing-indices", nil, samples, 2, false},
		{"index-count", []lumps.AmbientIndex{{AmbientSampleCount: 2, FirstAmbientSample: 0}}, samples, 2, false},
		{"span-past-end", []lumps.AmbientIndex{{AmbientSampleCount: 1, FirstAmbientSample: 2}, {AmbientSampleCount: 1, FirstAmbientSample: 0}}, samples, 2, false},
		{"wide-span-no-wrap", []lumps.AmbientIndex{{AmbientSampleCount: 65535, FirstAmbientSample: 65535}, {AmbientSampleCount: 1, FirstAmbientSample: 0}}, samples, 2, false},
		{"missing-samples", []lumps.AmbientIndex{{AmbientSampleCount: 1, FirstAmbientSample: 0}, {AmbientSampleCount: 1, FirstAmbientSample: 1}}, nil, 2, false},
		{"negative-leaves", nil, nil, -1, false},
		{"too-many-leaves", nil, nil, lumps.MaxMapLeafs + 1, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := ValidateAmbient(tc.indices, tc.samples, tc.leaves, context); (err == nil) != tc.valid {
				t.Fatalf("valid=%v err=%v", tc.valid, err)
			}
		})
	}
	// A zero-count reference can exceed sample count but remain a valid leaf ID.
	indices := []lumps.AmbientIndex{{AmbientSampleCount: 1, FirstAmbientSample: 0}, {AmbientSampleCount: 0, FirstAmbientSample: 2}, {AmbientSampleCount: 0, FirstAmbientSample: 0}}
	if err := ValidateAmbient(indices, samples[:1], 3, context); err != nil {
		t.Fatalf("leaf reference treated as sample span: %v", err)
	}
	maxIndices := make([]lumps.AmbientIndex, lumps.MaxMapLeafs)
	maxIndices[0].FirstAmbientSample = 65535
	if err := ValidateAmbient(maxIndices, nil, lumps.MaxMapLeafs, context); err != nil {
		t.Fatalf("rejected maximum representable leaf reference: %v", err)
	}
	legacy := ambientTestContext(true)
	if err := ValidateAmbient(nil, ambientExpectedRecords(false), 1, legacy); err == nil {
		t.Fatal("accepted extra legacy cube")
	}
	if err := ValidateAmbient(nil, ambientExpectedRecords(false)[:1], 2, legacy); err == nil {
		t.Fatal("accepted missing legacy cube")
	}
	if err := ValidateAmbient([]lumps.AmbientIndex{{}}, ambientExpectedRecords(false), 2, legacy); err == nil {
		t.Fatal("accepted ignored legacy indices")
	}
}

func TestAmbientIndexCountAndOutputBoundaries(t *testing.T) {
	context := ambientTestContext(false)
	context.LumpVersion = 0
	for _, count := range []int{0, lumps.MaxMapLeafs} {
		raw := make([]byte, count*4)
		data, err := DecodeAmbientIndices(raw, context)
		if err != nil || len(data) != count {
			t.Fatalf("index count %d: %v", count, err)
		}
		encoded, err := EncodeAmbientIndices(data, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("index count %d encode: %v", count, err)
		}
	}
	if _, err := DecodeAmbientIndices(make([]byte, (lumps.MaxMapLeafs+1)*4), context); err == nil {
		t.Fatal("decoded too many indices")
	}
	if _, err := EncodeAmbientIndices(make([]lumps.AmbientIndex, lumps.MaxMapLeafs+1), context); err == nil {
		t.Fatal("encoded too many indices")
	}
	for _, stride := range []int{4, 24, 28} {
		limit := (1<<31 - 1) / stride
		if size, err := ambientOutputSize(limit, stride); err != nil || size != limit*stride {
			t.Fatalf("representable output boundary %d: %d %v", stride, size, err)
		}
		if _, err := ambientOutputSize(limit+1, stride); err == nil {
			t.Fatalf("accepted output above signed32 boundary stride %d", stride)
		}
		if _, err := ambientOutputSize(int(^uint(0)>>1), stride); err == nil {
			t.Fatalf("accepted native-int maximum count stride %d", stride)
		}
	}
}

func TestAmbientRejectsWrongContexts(t *testing.T) {
	contexts := []format.Context{
		{}, {Profile: "sdk2013", BSPVersion: 21, LeafVersion: 1},
		{Profile: "sdk2013", BSPVersion: 20, LeafVersion: 1, ByteOrder: binary.BigEndian},
		{Profile: "future", BSPVersion: 20, LeafVersion: 1},
	}
	for _, version := range []int32{-1, 2, 20} {
		context := ambientTestContext(false)
		context.LeafVersion = version
		contexts = append(contexts, context)
	}
	for _, context := range contexts {
		if _, err := DecodeAmbientSamples(nil, context); err == nil {
			t.Fatalf("decoded samples with %+v", context)
		}
		if _, err := EncodeAmbientSamples(nil, context); err == nil {
			t.Fatalf("encoded samples with %+v", context)
		}
		if _, err := DecodeAmbientIndices(nil, context); err == nil {
			t.Fatalf("decoded indices with %+v", context)
		}
		if _, err := EncodeAmbientIndices(nil, context); err == nil {
			t.Fatalf("encoded indices with %+v", context)
		}
		if err := ValidateAmbient(nil, nil, 0, context); err == nil {
			t.Fatalf("validated %+v", context)
		}
	}
	for _, version := range []int32{-1, 2, 999} {
		context := ambientTestContext(false)
		context.LumpVersion = version
		if _, err := DecodeAmbientSamples(nil, context); err == nil {
			t.Fatal("decoded unsupported indexed sample revision")
		}
		if _, err := EncodeAmbientSamples(nil, context); err == nil {
			t.Fatal("encoded unsupported indexed sample revision")
		}
	}
	for _, version := range []int32{-1, 1, 2} {
		context := ambientTestContext(false)
		context.LumpVersion = version
		if _, err := DecodeAmbientIndices(nil, context); err == nil {
			t.Fatal("decoded unsupported index revision")
		}
		if _, err := EncodeAmbientIndices(nil, context); err == nil {
			t.Fatal("encoded unsupported index revision")
		}
	}
}

func TestReactiveDropAmbientSiblingLayoutsAndPadding(t *testing.T) {
	for _, legacy := range []bool{false, true} {
		context := format.Context{Profile: "reactivedrop", BSPVersion: 21, LeafVersion: 1, LegacyAmbient: legacy}
		versions := []int32{0, 1}
		if legacy {
			// The pinned consumer ignores garbage descriptors for old cubes;
			// only absence of the LDR index selects both channels' old width.
			versions = append(versions, -123, 999)
		}
		for _, version := range versions {
			context.LumpVersion = version
			raw, want := ambientWireRecords(!legacy), ambientExpectedRecords(!legacy)
			got, err := DecodeAmbientSamples(raw, context)
			if err != nil || !reflect.DeepEqual(got, want) {
				t.Fatalf("RD ambient legacy=%v v%d: got %+v err=%v want %+v", legacy, version, got, err, want)
			}
			encoded, err := EncodeAmbientSamples(want, context)
			if err != nil || !bytes.Equal(encoded, raw) {
				t.Fatalf("RD authored ambient legacy=%v v%d: %v", legacy, version, err)
			}
			got[1].Cube.Color[5].Exponent = -127
			mutated := append([]byte(nil), raw...)
			stride := 28
			if legacy {
				stride = 24
			}
			mutated[stride+23] = 0x81
			if !legacy {
				got[1].Pad = 0xef
				mutated[stride+27] = 0xef
			}
			encoded, err = EncodeAmbientSamples(got, context)
			if err != nil || !bytes.Equal(encoded, mutated) {
				t.Fatalf("RD ambient exponent/padding edit changed other fields: %v", err)
			}
			for n := 1; n < len(raw); n++ {
				if n == stride {
					continue
				}
				if records, err := DecodeAmbientSamples(raw[:n], context); err == nil || records != nil {
					t.Fatalf("RD ambient legacy=%v accepted partial length %d", legacy, n)
				}
			}
		}
	}
	context := format.Context{Profile: "reactivedrop", BSPVersion: 21, LeafVersion: 1, LumpVersion: 0}
	raw := []byte{0x03, 0x01, 0x02, 0x04, 0x05, 0x08, 0x07, 0x06}
	want := []lumps.AmbientIndex{
		{AmbientSampleCount: 0x0103, FirstAmbientSample: 0x0402},
		{AmbientSampleCount: 0x0805, FirstAmbientSample: 0x0607},
	}
	got, err := DecodeAmbientIndices(raw, context)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("RD ambient index offsets: got %+v err=%v", got, err)
	}
	encoded, err := EncodeAmbientIndices(want, context)
	if err != nil || !bytes.Equal(encoded, raw) {
		t.Fatalf("RD independent ambient index encoding: %v", err)
	}
	for n := 1; n < len(raw); n++ {
		if n == 4 {
			continue
		}
		if value, err := DecodeAmbientIndices(raw[:n], context); err == nil || value != nil {
			t.Fatalf("RD ambient indices accepted partial length %d", n)
		}
	}
	for _, leafVersion := range []int32{0, 1} {
		context.LeafVersion, context.LegacyAmbient, context.LumpVersion = leafVersion, true, 999
		if _, err := DecodeAmbientIndices(raw, context); err == nil {
			t.Fatal("RD old ambient layout consumed ignored separate indices")
		}
		if value, err := DecodeAmbientIndices(nil, context); err != nil || len(value) != 0 {
			t.Fatalf("RD old ambient empty indices: %v", err)
		}
		if raw, err := EncodeAmbientIndices(want, context); err == nil || raw != nil {
			t.Fatal("RD old ambient layout wrote ignored separate indices")
		}
	}
	context.LeafVersion, context.LegacyAmbient = 0, false
	if value, err := DecodeAmbientSamples(nil, context); err != nil || len(value) != 0 {
		t.Fatalf("RD leaf-v0 unused ambient channel: %v", err)
	}
	if _, err := DecodeAmbientSamples(ambientWireRecords(false), context); err == nil {
		t.Fatal("RD leaf-v0 consumed separate lighting instead of inline cube")
	}
	if raw, err := EncodeAmbientSamples(ambientExpectedRecords(false), context); err == nil || raw != nil {
		t.Fatal("RD leaf-v0 wrote separate lighting")
	}
}

func TestReactiveDropAmbientRejectsWrongContexts(t *testing.T) {
	for _, context := range []format.Context{
		{Profile: "reactivedrop", BSPVersion: 19, LeafVersion: 1},
		{Profile: "reactivedrop", BSPVersion: 20, LeafVersion: 1},
		{Profile: "reactivedrop", BSPVersion: 22, LeafVersion: 1},
		{Profile: "reactivedrop", BSPVersion: 21, LeafVersion: -1},
		{Profile: "reactivedrop", BSPVersion: 21, LeafVersion: 2},
		{Profile: "reactivedrop", BSPVersion: 21, LeafVersion: 1, LumpVersion: -1},
		{Profile: "reactivedrop", BSPVersion: 21, LeafVersion: 1, LumpVersion: 2},
		{Profile: "reactivedrop", BSPVersion: 21, LeafVersion: 1, ByteOrder: binary.BigEndian},
		{Profile: "reactivedrop", BSPVersion: 21, LeafVersion: 1, HeaderLayout: format.L4D2Header},
	} {
		if _, err := DecodeAmbientSamples(ambientWireRecords(true), context); err == nil {
			t.Fatalf("RD ambient samples accepted context %+v", context)
		}
		if raw, err := EncodeAmbientSamples(ambientExpectedRecords(true), context); err == nil || raw != nil {
			t.Fatalf("RD ambient samples encoded context %+v", context)
		}
		if _, err := DecodeAmbientIndices([]byte{1, 0, 0, 0}, context); err == nil {
			t.Fatalf("RD ambient indices accepted context %+v", context)
		}
		if raw, err := EncodeAmbientIndices([]lumps.AmbientIndex{{AmbientSampleCount: 1}}, context); err == nil || raw != nil {
			t.Fatalf("RD ambient indices encoded context %+v", context)
		}
	}
	// Sample revision 1 is valid, but index revision 1 is not.
	context := format.Context{Profile: "reactivedrop", BSPVersion: 21, LeafVersion: 1, LumpVersion: 1}
	if _, err := DecodeAmbientIndices([]byte{1, 0, 0, 0}, context); err == nil {
		t.Fatal("RD ambient index accepted sample revision")
	}
	if raw, err := EncodeAmbientIndices([]lumps.AmbientIndex{{AmbientSampleCount: 1}}, context); err == nil || raw != nil {
		t.Fatal("RD ambient index encoded sample revision")
	}
}

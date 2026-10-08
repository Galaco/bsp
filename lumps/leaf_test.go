package lumps

import (
	"bytes"
	"testing"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/primitives/common"
	primitives "github.com/galaco/bsp/primitives/leaf"
)

// These wire records are authored independently of the encoder. Each contains
// distinct signed fields and bounds so a wrong stride cannot decode the second.
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
	first = append(first, 0, 0)
	second = append(second, 0, 0)
	return append(first, second...)
}

func leafExpectedRecords(version int32) []primitives.Leaf {
	data := []primitives.Leaf{
		{
			Contents: 1, Cluster: 2, BitField: 0x5c3d,
			Mins: [3]int16{3, 7, 14}, Maxs: [3]int16{13, 36, 51},
			FirstLeafFace: 32, NumLeafFaces: 11, FirstLeafBrush: 21,
			NumLeafBrushes: 76, LeafWaterDataID: 621,
		},
		{
			Contents: -2, Cluster: -1, BitField: -32253,
			Mins: [3]int16{-1, -2, -3}, Maxs: [3]int16{4, 5, 6},
			FirstLeafFace: 0x1234, NumLeafFaces: 0x5678, FirstLeafBrush: 0x9abc,
			NumLeafBrushes: 0xdef0, LeafWaterDataID: -1,
		},
	}
	if version == 0 {
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

func TestLeafWireRoundTrip(t *testing.T) {
	for _, tc := range []struct {
		name       string
		version    int32
		bspVersion int32
		stride     int
	}{
		{"v0-in-bsp21", 0, 21, 56},
		{"v1-in-bsp19", 1, 19, 32},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var lump Leaf
			lump.SetContext(format.Context{BSPVersion: tc.bspVersion, LumpVersion: tc.version})
			raw := leafWireRecords(tc.version)
			if err := lump.Unmarshall(raw); err != nil {
				t.Fatal(err)
			}
			expected := leafExpectedRecords(tc.version)
			if len(lump.data) != len(expected) {
				t.Fatalf("decoded %d records, want %d", len(lump.data), len(expected))
			}
			for i := range expected {
				if lump.data[i] != expected[i] {
					t.Errorf("record %d: got %+v, want %+v", i, lump.data[i], expected[i])
				}
			}
			encoded, err := lump.Marshall()
			if err != nil {
				t.Fatal(err)
			}
			if len(encoded) != 2*tc.stride || !bytes.Equal(encoded, raw) {
				t.Errorf("wire mismatch: got %x, want %x", encoded, raw)
			}
			// Encoding independently supplied records must use the same layout.
			lump.data = expected
			encoded, err = lump.Marshall()
			if err != nil || !bytes.Equal(encoded, raw) {
				t.Fatalf("authored records encoded as %x, err %v; want %x", encoded, err, raw)
			}
		})
	}
}

func TestLeafRejectsMalformedLength(t *testing.T) {
	for _, version := range []int32{0, 1} {
		raw := leafWireRecords(version)
		for _, malformed := range [][]byte{raw[:len(raw)-1], append(raw, 1), {1}} {
			var lump Leaf
			lump.SetContext(format.Context{LumpVersion: version})
			if err := lump.Unmarshall(malformed); err == nil {
				t.Errorf("version %d accepted malformed length %d", version, len(malformed))
			}
		}
	}
}

func TestLeafCountBoundaries(t *testing.T) {
	for _, version := range []int32{0, 1} {
		stride := 32
		if version == 0 {
			stride = 56
		}
		var lump Leaf
		lump.SetContext(format.Context{LumpVersion: version})
		for _, count := range []int{0, MaxMapLeafs} {
			raw := make([]byte, count*stride)
			if err := lump.Unmarshall(raw); err != nil {
				t.Fatalf("version %d rejected count %d: %v", version, count, err)
			}
			encoded, err := lump.Marshall()
			if err != nil || !bytes.Equal(encoded, raw) {
				t.Fatalf("version %d failed round trip for count %d: %v", version, count, err)
			}
		}
		if err := lump.Unmarshall(make([]byte, (MaxMapLeafs+1)*stride)); err == nil {
			t.Errorf("version %d accepted too many records", version)
		}
		lump.data = make([]primitives.Leaf, MaxMapLeafs+1)
		if _, err := lump.Marshall(); err == nil {
			t.Errorf("version %d encoded too many records", version)
		}
	}
}

func TestLeafRejectsUnsupportedVersion(t *testing.T) {
	for _, version := range []int32{-1, 2, 19, 20, 21} {
		var lump Leaf
		lump.SetContext(format.Context{LumpVersion: version})
		if err := lump.Unmarshall(nil); err == nil {
			t.Errorf("decoded unsupported version %d", version)
		}
		if _, err := lump.Marshall(); err == nil {
			t.Errorf("encoded unsupported version %d", version)
		}
	}
}

func TestLeafIgnoresNonsemanticPadding(t *testing.T) {
	for _, version := range []int32{0, 1} {
		raw := leafWireRecords(version)
		stride := len(raw) / 2
		for i := 0; i < 2; i++ {
			raw[(i+1)*stride-2] = 0xab
			raw[(i+1)*stride-1] = 0xcd
		}
		var lump Leaf
		lump.SetContext(format.Context{LumpVersion: version})
		if err := lump.Unmarshall(raw); err != nil {
			t.Fatal(err)
		}
		expected := leafExpectedRecords(version)
		for i := range expected {
			if lump.data[i] != expected[i] {
				t.Fatalf("padding changed record %d: %+v", i, lump.data[i])
			}
		}
	}
}

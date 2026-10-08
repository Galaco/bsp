package lumps_test

import (
	"testing"

	"github.com/galaco/bsp/lumps"
	facevalue "github.com/galaco/bsp/primitives/face"
)

func TestLeafUnsignedPackedFields(t *testing.T) {
	for _, tc := range []struct{ word, area, flags uint16 }{{0, 0, 0}, {0x007f, 127, 0}, {0x0080, 128, 0}, {0x01ff, 511, 0}, {0x0200, 0, 1}, {0x81ff, 511, 64}, {0xffff, 511, 127}} {
		leaf := lumps.Leaf{BitField: tc.word}
		if leaf.Area() != tc.area || leaf.Flags() != tc.flags {
			t.Fatalf("word%04x area%d flags%d", tc.word, leaf.Area(), leaf.Flags())
		}
	}
	for _, tc := range []struct{ input, word uint16 }{{0, 0x8000}, {127, 0x807f}, {128, 0x8080}, {511, 0x81ff}, {512, 0x8000}, {0xffff, 0x81ff}} {
		leaf := lumps.Leaf{BitField: 0x8080}
		leaf.SetArea(tc.input)
		if leaf.BitField != tc.word || leaf.Flags() != 64 {
			t.Fatalf("SetArea(%d) word%04x flags%d", tc.input, leaf.BitField, leaf.Flags())
		}
	}
	for _, tc := range []struct{ input, word uint16 }{{0, 0x01ff}, {1, 0x03ff}, {64, 0x81ff}, {127, 0xffff}, {128, 0x01ff}, {0xffff, 0xffff}} {
		leaf := lumps.Leaf{BitField: 0x81ff}
		leaf.SetFlags(tc.input)
		if leaf.BitField != tc.word || leaf.Area() != 511 {
			t.Fatalf("SetFlags(%d) word%04x area%d", tc.input, leaf.BitField, leaf.Area())
		}
	}
}

func TestFaceCountAndInvertedShadowStateTransitions(t *testing.T) {
	for _, count := range []uint16{0, 1, 32767} {
		face := facevalue.Face{NumPrims: count}
		face.SetDynamicShadowsEnabled(false)
		if face.NumPrims != count|0x8000 || face.PrimitiveCount() != count || face.DynamicShadowsEnabled() {
			t.Fatalf("disable count%d word%04x", count, face.NumPrims)
		}
		face.SetDynamicShadowsEnabled(true)
		if face.NumPrims != count || face.PrimitiveCount() != count || !face.DynamicShadowsEnabled() {
			t.Fatalf("enable count%d word%04x", count, face.NumPrims)
		}
	}
	for _, tc := range []struct{ initial, input, word uint16 }{{0x8000, 1, 0x8001}, {0x8001, 32767, 0xffff}, {0xffff, 0, 0x8000}, {0xffff, 32768, 0x8000}, {0x0001, 0xffff, 0x7fff}} {
		face := facevalue.Face{NumPrims: tc.initial}
		enabled := face.DynamicShadowsEnabled()
		face.SetPrimitiveCount(tc.input)
		if face.NumPrims != tc.word || face.DynamicShadowsEnabled() != enabled {
			t.Fatalf("count input%d initial%04x word%04x", tc.input, tc.initial, face.NumPrims)
		}
	}
}

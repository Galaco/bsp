package lumps_test

import (
	"testing"

	"github.com/galaco/bsp/lumps"
)

func TestOverlayPackedFields(t *testing.T) {
	for _, tc := range []struct {
		word        uint16
		faceCount   uint16
		renderOrder uint16
	}{
		{0x0000, 0, 0},
		{0x0040, 64, 0},
		{0x3fff, 0x3fff, 0},
		{0x4000, 0, 1},
		{0xc000, 0, 3},
		{0x4040, 64, 1},
		{0xc3ff, 0x3ff, 3},
		{0xffff, 0x3fff, 3},
	} {
		t.Run("read", func(t *testing.T) {
			overlay := lumps.Overlay{FaceCountAndRenderOrder: tc.word}
			if overlay.FaceCount() != tc.faceCount {
				t.Fatalf("word=%04x faceCount=%d want=%d", tc.word, overlay.FaceCount(), tc.faceCount)
			}
			if overlay.RenderOrder() != tc.renderOrder {
				t.Fatalf("word=%04x renderOrder=%d want=%d", tc.word, overlay.RenderOrder(), tc.renderOrder)
			}
		})
	}

	for _, tc := range []struct {
		input uint16
		word  uint16
	}{
		{0, 0x8000},
		{32, 0x8020},
		{64, 0x8040},
		{0x3fff, 0xbfff},
		{65, 0x8000},
		{0xffff, 0xbfff},
	} {
		t.Run("setFaceCount", func(t *testing.T) {
			overlay := lumps.Overlay{FaceCountAndRenderOrder: 0x8000}
			err := overlay.SetFaceCount(tc.input)
			if tc.input > 64 {
				if err == nil {
					t.Fatalf("SetFaceCount(%d) expected error", tc.input)
				}
				if overlay.FaceCountAndRenderOrder != 0x8000 {
					t.Fatalf("SetFaceCount(%d) changed word to %04x on error", tc.input, overlay.FaceCountAndRenderOrder)
				}
				return
			}
			if err != nil {
				t.Fatalf("SetFaceCount(%d) unexpected error: %v", tc.input, err)
			}
			if overlay.FaceCountAndRenderOrder != tc.word {
				t.Fatalf("SetFaceCount(%d) word=%04x want=%04x", tc.input, overlay.FaceCountAndRenderOrder, tc.word)
			}
			if overlay.RenderOrder() != 2 {
				t.Fatalf("SetFaceCount(%d) clobbered renderOrder=%d want=2", tc.input, overlay.RenderOrder())
			}
		})
	}

	for _, tc := range []struct {
		input uint16
		word  uint16
	}{
		{0, 0x0040},
		{1, 0x4040},
		{2, 0x8040},
		{3, 0xc040},
		{4, 0x0040},
		{0xffff, 0xc040},
	} {
		t.Run("setRenderOrder", func(t *testing.T) {
			overlay := lumps.Overlay{FaceCountAndRenderOrder: 0x0040}
			err := overlay.SetRenderOrder(tc.input)
			if tc.input > 3 {
				if err == nil {
					t.Fatalf("SetRenderOrder(%d) expected error", tc.input)
				}
				if overlay.FaceCountAndRenderOrder != 0x0040 {
					t.Fatalf("SetRenderOrder(%d) changed word to %04x on error", tc.input, overlay.FaceCountAndRenderOrder)
				}
				return
			}
			if err != nil {
				t.Fatalf("SetRenderOrder(%d) unexpected error: %v", tc.input, err)
			}
			if overlay.FaceCountAndRenderOrder != tc.word {
				t.Fatalf("SetRenderOrder(%d) word=%04x want=%04x", tc.input, overlay.FaceCountAndRenderOrder, tc.word)
			}
			if overlay.FaceCount() != 64 {
				t.Fatalf("SetRenderOrder(%d) clobbered faceCount=%d want=64", tc.input, overlay.FaceCount())
			}
		})
	}
}

func TestWaterOverlayPackedFields(t *testing.T) {
	for _, tc := range []struct {
		word        uint16
		faceCount   uint16
		renderOrder uint16
	}{
		{0x0000, 0, 0},
		{0x0100, 256, 0},
		{0x3fff, 0x3fff, 0},
		{0x4000, 0, 1},
		{0xc000, 0, 3},
		{0x4100, 256, 1},
		{0xffff, 0x3fff, 3},
	} {
		t.Run("read", func(t *testing.T) {
			overlay := lumps.WaterOverlay{FaceCountAndRenderOrder: tc.word}
			if overlay.FaceCount() != tc.faceCount {
				t.Fatalf("word=%04x faceCount=%d want=%d", tc.word, overlay.FaceCount(), tc.faceCount)
			}
			if overlay.RenderOrder() != tc.renderOrder {
				t.Fatalf("word=%04x renderOrder=%d want=%d", tc.word, overlay.RenderOrder(), tc.renderOrder)
			}
		})
	}

	for _, tc := range []struct {
		input uint16
		word  uint16
	}{
		{0, 0x8000},
		{128, 0x8080},
		{256, 0x8100},
		{0x3fff, 0xbfff},
		{257, 0x8000},
		{0xffff, 0xbfff},
	} {
		t.Run("setFaceCount", func(t *testing.T) {
			overlay := lumps.WaterOverlay{FaceCountAndRenderOrder: 0x8000}
			err := overlay.SetFaceCount(tc.input)
			if tc.input > 256 {
				if err == nil {
					t.Fatalf("SetFaceCount(%d) expected error", tc.input)
				}
				if overlay.FaceCountAndRenderOrder != 0x8000 {
					t.Fatalf("SetFaceCount(%d) changed word to %04x on error", tc.input, overlay.FaceCountAndRenderOrder)
				}
				return
			}
			if err != nil {
				t.Fatalf("SetFaceCount(%d) unexpected error: %v", tc.input, err)
			}
			if overlay.FaceCountAndRenderOrder != tc.word {
				t.Fatalf("SetFaceCount(%d) word=%04x want=%04x", tc.input, overlay.FaceCountAndRenderOrder, tc.word)
			}
			if overlay.RenderOrder() != 2 {
				t.Fatalf("SetFaceCount(%d) clobbered renderOrder=%d want=2", tc.input, overlay.RenderOrder())
			}
		})
	}

	for _, tc := range []struct {
		input uint16
		word  uint16
	}{
		{0, 0x0100},
		{1, 0x4100},
		{2, 0x8100},
		{3, 0xc100},
		{4, 0x0100},
		{0xffff, 0xc100},
	} {
		t.Run("setRenderOrder", func(t *testing.T) {
			overlay := lumps.WaterOverlay{FaceCountAndRenderOrder: 0x0100}
			err := overlay.SetRenderOrder(tc.input)
			if tc.input > 3 {
				if err == nil {
					t.Fatalf("SetRenderOrder(%d) expected error", tc.input)
				}
				if overlay.FaceCountAndRenderOrder != 0x0100 {
					t.Fatalf("SetRenderOrder(%d) changed word to %04x on error", tc.input, overlay.FaceCountAndRenderOrder)
				}
				return
			}
			if err != nil {
				t.Fatalf("SetRenderOrder(%d) unexpected error: %v", tc.input, err)
			}
			if overlay.FaceCountAndRenderOrder != tc.word {
				t.Fatalf("SetRenderOrder(%d) word=%04x want=%04x", tc.input, overlay.FaceCountAndRenderOrder, tc.word)
			}
			if overlay.FaceCount() != 256 {
				t.Fatalf("SetRenderOrder(%d) clobbered faceCount=%d want=256", tc.input, overlay.FaceCount())
			}
		})
	}
}

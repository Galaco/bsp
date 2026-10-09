package lumps

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

// WaterOverlay is the complete revision-0, 1120-byte lump 50 record.
// Inactive Faces slots are retained wire data, not spare storage to clear.
type WaterOverlay struct {
	ID                      int32
	TexInfo                 int16
	FaceCountAndRenderOrder uint16
	Faces                   [256]int32
	U, V                    [2]float32
	UVPoints                [4]mgl32.Vec3
	Origin, BasisNormal     mgl32.Vec3
}

func (record *WaterOverlay) FaceCount() uint16   { return record.FaceCountAndRenderOrder & 0x3fff }
func (record *WaterOverlay) RenderOrder() uint16 { return record.FaceCountAndRenderOrder >> 14 }

// SetFaceCount preserves the render order and all inactive face slots.
// Invalid values leave the record unchanged.
func (record *WaterOverlay) SetFaceCount(count uint16) error {
	if count > 256 {
		return fmt.Errorf("water overlay face count %d exceeds 256", count)
	}
	record.FaceCountAndRenderOrder = record.FaceCountAndRenderOrder&0xc000 | count
	return nil
}

// SetRenderOrder preserves the face count; invalid values change nothing.
func (record *WaterOverlay) SetRenderOrder(order uint16) error {
	if order > 3 {
		return fmt.Errorf("water overlay render order %d exceeds 3", order)
	}
	record.FaceCountAndRenderOrder = record.FaceCountAndRenderOrder&0x3fff | order<<14
	return nil
}

// OverlaySystemLevel is the four-byte modern lump 61 record. These fields
// retain the branch's byte codes without imposing an invented level ordering.
type OverlaySystemLevel struct {
	MinCPULevel, MaxCPULevel, MinGPULevel, MaxGPULevel uint8
}

// DisplacementMultiblend is the 80-byte modern lump 63 record. Every float
// component, including the four per-channel colors, is retained verbatim.
type DisplacementMultiblend struct {
	Weights, Alpha mgl32.Vec4
	Colors         [4]mgl32.Vec3
}

// PropCollision retains Reactive Drop's complete eight-byte lump 22 record.
// Public wire sources do not establish its model-dictionary association.
type PropCollision struct {
	HullCount, HullStart int32
}

// PropHull retains all sixteen bytes of Reactive Drop's lump 23 record.
// SurfaceProperty and span scalars keep their original signed wire values.
type PropHull struct {
	VertexCount, VertexStart, SurfaceProperty int32
	Contents                                  uint32
}

// PropHullTriangles is Reactive Drop's eight-byte lump 25 record. IndexStart
// and IndexCount retain the signed wire span; the companion blob stays opaque.
type PropHullTriangles struct {
	IndexStart, IndexCount int32
}

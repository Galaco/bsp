package lumps

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

// Overlay is the complete revision-0, 352-byte lump 45 record. Inactive Faces
// slots remain wire data and must not be cleared when editing active faces.
type Overlay struct {
	ID                      int32
	TexInfo                 int16
	FaceCountAndRenderOrder uint16
	Faces                   [64]int32
	U, V                    [2]float32
	UVPoints                [4]mgl32.Vec3
	Origin, BasisNormal     mgl32.Vec3
}

// FaceCount returns the unsigned lower fourteen packed bits.
func (record *Overlay) FaceCount() uint16 { return record.FaceCountAndRenderOrder & 0x3fff }

// RenderOrder returns the unsigned upper two packed bits.
func (record *Overlay) RenderOrder() uint16 { return record.FaceCountAndRenderOrder >> 14 }

// SetFaceCount retains the render order and leaves the record unchanged on error.
func (record *Overlay) SetFaceCount(count uint16) error {
	if count > 64 {
		return fmt.Errorf("overlay face count %d exceeds 64", count)
	}
	record.FaceCountAndRenderOrder = record.FaceCountAndRenderOrder&0xc000 | count
	return nil
}

// SetRenderOrder retains the face count and leaves the record unchanged on error.
func (record *Overlay) SetRenderOrder(order uint16) error {
	if order > 3 {
		return fmt.Errorf("overlay render order %d exceeds 3", order)
	}
	record.FaceCountAndRenderOrder = record.FaceCountAndRenderOrder&0x3fff | order<<14
	return nil
}

// Cubemap is a revision-0, 16-byte lump 42 sample. Origin is snapped integer
// coordinates; Size is the original wire size code, not an expanded dimension.
type Cubemap struct {
	Origin  [3]int32
	Size    uint8
	Padding [3]byte
}

// Primitive is a revision-0, ten-byte lump 37 record. Type 0 denotes a triangle
// list and Type 1 a triangle strip. Padding is preserved, not normalized.
type Primitive struct {
	Type                                         uint8
	Padding                                      uint8
	FirstIndex, IndexCount, FirstVert, VertCount uint16
}

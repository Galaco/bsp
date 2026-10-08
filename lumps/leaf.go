package lumps

import "github.com/galaco/bsp/primitives/common"

// MaxMapLeafs is the maximum number of leaves supported by the engine.
const MaxMapLeafs = 65536

const (
	LeafFlagsSky    = 0x01
	LeafFlagsRadial = 0x02
	LeafFlagsSky2D  = 0x04
)

// Leaf is a convex volume from lump 10. The descriptor version selects whether
// LightSample is stored inline; Padding preserves the two trailing wire bytes.
type Leaf struct {
	Contents          int32
	Cluster           int16
	BitField          uint16
	Mins              [3]int16
	Maxs              [3]int16
	FirstLeafFace     uint16
	NumLeafFaces      uint16
	FirstLeafBrush    uint16
	NumLeafBrushes    uint16
	LeafWaterDataID   int16
	LightSample       common.CompressedLightCube
	HasInlineLighting bool
	Padding           [2]byte
}

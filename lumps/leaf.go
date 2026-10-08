package lumps

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
	LightSample       CompressedLightCube
	HasInlineLighting bool
	Padding           [2]byte
}

// Area is the unsigned nine-bit area stored in BitField.
func (leaf *Leaf) Area() uint16 { return leaf.BitField & 0x01ff }

// Flags is the unsigned seven-bit flag value, never an arithmetic signed shift.
func (leaf *Leaf) Flags() uint16 { return leaf.BitField >> 9 }

// SetArea masks to nine bits and preserves every flag bit.
func (leaf *Leaf) SetArea(area uint16) {
	leaf.BitField = leaf.BitField&0xfe00 | area&0x01ff
}

// SetFlags masks to seven bits and preserves every area bit.
func (leaf *Leaf) SetFlags(flags uint16) {
	leaf.BitField = leaf.BitField&0x01ff | (flags&0x007f)<<9
}

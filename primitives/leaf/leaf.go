package leaf

const BITMASK_LOWER9 = 0x1FF // 511 (2^9 - 1)
const BITMASK_LOWER7 = 0x7F // 127 (2^7 - 1)

// NOTE: Only 7-bits stored!!!
const LEAF_FLAGS_SKY	= 0x01		// This leaf has 3D sky in its PVS
const LEAF_FLAGS_RADIAL	= 0x02		// This leaf culled away some portals due to radial vis
const LEAF_FLAGS_SKY2D	= 0x04		// This leaf has 2D sky in its PVS

type Leaf struct {
	Contents int32
	Cluster int16
	BitField int16 //C Union of char Name || Area:9 && Flags:7
	Mins [3]int16
	Maxs [3]int16
	FirstLeafFace uint16
	NumLeafFaces uint16
	FirstLeafBrush uint16
	NumLeafBrushes uint16
	LeafWaterDataID int16
	_ [2]byte
}

// Get area (first 9 bits)
func (b *Leaf) Area() int16{
	return int16((b.BitField >> 9) & BITMASK_LOWER9)
}
// Set area (first 9 bits)
func (b *Leaf) SetArea(area int16) {
	v := b.BitField
	b.BitField = int16((v & BITMASK_LOWER9) | (area))
}

// Get flags (second 7 bits)
func (b *Leaf) Flags() int16{
	return int16((b.BitField) & BITMASK_LOWER7)
}
// Set flags (second 7 bits)
func (b *Leaf) SetFlags(flags int16) {
	v := b.BitField
	b.BitField = int16((v & BITMASK_LOWER7) | (int16(flags) << 9))
}
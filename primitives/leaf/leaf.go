package leaf

const bitmaskLower9 = 0x1FF // 511 (2^9 - 1)
const bitmaskLower7 = 0x7F  // 127 (2^7 - 1)

// Leaf represents a single convex volume that contains 0 or more
// faces (or entities) within its bounds
type Leaf struct {
	// Contents
	Contents int32
	// Cluster that this leaf is a part of
	Cluster int16
	// BitField is a C Union of char Name || Area:9 && Flags:7
	BitField int16
	// Mins is this leafs bounding volumes minimum
	Mins [3]int16
	// Maxs  is this leafs bounding volumes maximum
	Maxs [3]int16
	// FirstLeafFace index into LeafFaces lump data
	FirstLeafFace uint16
	// NumLeafFaces is number of LeafFaces in this Leaf
	NumLeafFaces uint16
	// FirstLeafBrush is index into LeafBrushes lump data
	FirstLeafBrush uint16
	// NumLeafBrushes is number of LeafBrushes in this Leaf
	NumLeafBrushes uint16
	// LeafWaterDataID
	LeafWaterDataID int16
	_               [2]byte
}

// Area returns area (first 9 bits)
func (b *Leaf) Area() int16 {
	return int16((b.BitField) & bitmaskLower7)
}

// SetArea sets area (first 9 bits)
func (b *Leaf) SetArea(area int16) {
	v := b.BitField
	b.BitField = int16((v & bitmaskLower9) | (area))
}

// Flags returns flags (second 7 bits)
func (b *Leaf) Flags() int16 {
	return int16((b.BitField >> 9) & bitmaskLower9)
}

// SetFlags sets flags (second 7 bits)
func (b *Leaf) SetFlags(flags int16) {
	v := b.BitField
	b.BitField = int16((v & bitmaskLower7) | (int16(flags) << 9))
}

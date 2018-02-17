package leaf
type Leaf struct {
	Contents int32
	Cluster int16
	Area int16 // NOTE: Actually first 9 bits of a short, but not implemented
	Flags int16 // NOTE: Actually second 7 bits of a short, but not implemented
	Mins [3]int16
	Maxs [3]int16
	FirstLeafFace uint16
	NumLeafFaces uint16
	FirstLeafBrush uint16
	NumLeafBrushes uint16
	LeafWaterDataID int16
}

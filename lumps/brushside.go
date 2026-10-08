package lumps

// BrushSide is an eight-byte brush side. SDK2013 and L4D2 store Bevel as a
// signed short; CS:GO stores byte-sized Bevel and Thin. HasThin records that
// dialect explicitly, even when Thin is zero.
type BrushSide struct {
	PlaneNum uint16
	TexInfo  int16
	DispInfo int16
	Bevel    int16
	Thin     uint8
	HasThin  bool
}

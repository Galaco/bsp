package brushside

// BrushSide
type BrushSide struct {
	// PlaneNum is index into planes
	PlaneNum uint16
	// TexInfo is index into TexInfo slice
	TexInfo int16
	// DispInfo is index into DispInfo slice
	DispInfo int16
	// Bevel
	Bevel int16
}

package brushside

// BrushSide
type BrushSide struct {
	// PlaneNum is index into planes
	PlaneNum uint16 `json:"planeNum"`
	// TexInfo is index into TexInfo slice
	TexInfo int16 `json:"texInfo"`
	// DispInfo is index into DispInfo slice
	DispInfo int16 `json:"dispInfo"`
	// Bevel
	Bevel int16 `json:"bevel"`
}

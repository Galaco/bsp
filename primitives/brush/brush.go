package brush

// Brush
type Brush struct {
	// FirstSide index of first side of a brush
	FirstSide int32 `json:"firstSide"`
	// NumSides is number of sides this brush has
	NumSides int32 `json:"numSides"`
	// Contents
	Contents int32 `json:"contents"`
}

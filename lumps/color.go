package lumps

// ColorRGBExponent32 is a 4 byte color representation (RGBExp)
type ColorRGBExponent32 struct {
	// R red
	R uint8
	// G green
	G uint8
	// B blue
	B uint8
	// Exponent exponent
	Exponent int8
}

// CompressedLightCube
type CompressedLightCube struct {
	// Color
	Color [6]ColorRGBExponent32
}

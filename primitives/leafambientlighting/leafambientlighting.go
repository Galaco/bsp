package leafambientlighting

import primitives "github.com/galaco/bsp/primitives/common"

// LeafAmbientLighting
type LeafAmbientLighting struct {
	// Cube
	Cube CompressedLightCube
	// X x
	X byte
	// Y y
	Y byte
	// Z z
	Z byte
	// Pad is padding to 4 bytes (any other purpose unknown)
	Pad byte
}

// CompressedLightCube
type CompressedLightCube struct {
	// Color
	Color [6]primitives.ColorRGBExponent32
}

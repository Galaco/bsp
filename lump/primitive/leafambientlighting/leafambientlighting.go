package leafambientlighting

import (
	primitives "github.com/galaco/bsp/lump/primitive/common"
)

// LeafAmbientLighting
type LeafAmbientLighting struct {
	// Cube
	Cube primitives.CompressedLightCube `json:"cube"`
	// X x
	X byte `json:"x"`
	// Y y
	Y byte `json:"y"`
	// Z z
	Z byte `json:"z"`
	// Pad is padding to 4 bytes (any other purpose unknown)
	_ byte
}

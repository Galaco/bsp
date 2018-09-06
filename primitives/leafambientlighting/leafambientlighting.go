package leafambientlighting

import primitives "github.com/galaco/bsp/primitives/common"

type LeafAmbientLighting struct {
	Cube CompressedLightCube
	X    byte
	Y    byte
	Z    byte
	Pad  byte
}

type CompressedLightCube struct {
	Color [6]primitives.ColorRGBExponent32
}

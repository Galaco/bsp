package leafambientlighting

import "github.com/galaco/bsp/lumps/datatypes/common"

type LeafAmbientLighting struct {
	Cube CompressedLightCube
	X byte
	Y byte
	Z byte
	Pad byte
}

type CompressedLightCube struct {
	Color [6]common.ColorRGBExponent32
}
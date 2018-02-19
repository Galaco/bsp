package leafambientlighting

type LeafAmbientLighting struct {
	Cube CompressedLightCube
	X byte
	Y byte
	Z byte
	Pad byte
}

type CompressedLightCube struct {
	Color [6]ColorRGBExponent32
}

type ColorRGBExponent32 struct {
	R byte
	G byte
	B byte
	Exponent byte
}
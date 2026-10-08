package lumps

import "github.com/galaco/bsp/primitives/common"

// AmbientSample is an ambient light cube from lump 55 or 56. Legacy records
// store only Cube; indexed records also store position and the explicit pad byte.
// HasPosition identifies the selected wire layout, not whether XYZ is nonzero.
type AmbientSample struct {
	Cube        common.CompressedLightCube
	X           byte
	Y           byte
	Z           byte
	Pad         byte
	HasPosition bool
}

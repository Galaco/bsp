package leaf

// NOTE: Only 7-bits stored!!!
const (
	// LeafFlagsSky This leaf has 3D sky in its PVS
	LeafFlagsSky = 0x01
	// LeafFlagsRadial This leaf culled away some portals due to radial vis
	LeafFlagsRadial = 0x02
	// LeafFlagsSky2D This leaf has 2D sky in its PVS
	LeafFlagsSky2D = 0x04
)

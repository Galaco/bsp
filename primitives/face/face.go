package face

// Face
type Face struct {
	// Planenum
	Planenum uint16
	// Side
	Side byte
	// OnNode
	OnNode byte
	// FirstEdge is index into the edges lump data
	FirstEdge int32
	// NumEdges is the number of edges that make up this face
	NumEdges int16
	// TexInfo is index into texinfos lump
	TexInfo int16
	// DispInfo is index into dispinfos lump
	DispInfo int16
	// SurfaceFogVolumeID
	SurfaceFogVolumeID int16
	// Styles
	Styles [4]byte
	// Lightofs is offset into lighting lump (-1 if unlit)
	Lightofs int32
	// Area is total size of this side (units^2)
	Area float32
	// LightmapTextureMinsInLuxels
	LightmapTextureMinsInLuxels [2]int32
	// LightmapTextureSizeInLuxels is XY size of lightmap data
	LightmapTextureSizeInLuxels [2]int32
	// OrigFace
	OrigFace int32
	// NumPrims retains the low15-bit count and high-bit dynamic-shadow disable.
	NumPrims uint16
	// FirstPrimID
	FirstPrimID uint16
	// SmoothingGroups
	SmoothingGroups uint32
}

// PrimitiveCount excludes the dynamic-shadow disable bit.
func (face *Face) PrimitiveCount() uint16 { return face.NumPrims & 0x7fff }

// SetPrimitiveCount masks to fifteen bits and preserves shadow state.
func (face *Face) SetPrimitiveCount(count uint16) {
	face.NumPrims = face.NumPrims&0x8000 | count&0x7fff
}

// DynamicShadowsEnabled reports the inverted high-bit convention.
func (face *Face) DynamicShadowsEnabled() bool { return face.NumPrims&0x8000 == 0 }

// SetDynamicShadowsEnabled changes only the high disable bit.
func (face *Face) SetDynamicShadowsEnabled(enabled bool) {
	if enabled {
		face.NumPrims &= 0x7fff
	} else {
		face.NumPrims |= 0x8000
	}
}

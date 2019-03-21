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
	// NumPrims
	NumPrims uint16
	// FirstPrimID
	FirstPrimID uint16
	// SmoothingGroups
	SmoothingGroups uint32
}

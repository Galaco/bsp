package face

// Face
type Face struct {
	// Planenum
	Planenum uint16 `json:"planenum"`
	// Side
	Side byte `json:"side"`
	// OnNode
	OnNode byte `json:"onNode"`
	// FirstEdge is index into the edges lump data
	FirstEdge int32 `json:"firstEdge"`
	// NumEdges is the number of edges that make up this face
	NumEdges int16 `json:"numEdges"`
	// TexInfo is index into texinfos lump
	TexInfo int16 `json:"texInfo"`
	// DispInfo is index into dispinfos lump
	DispInfo int16 `json:"dispInfo"`
	// SurfaceFogVolumeID
	SurfaceFogVolumeID int16 `json:"surfaceFogVolumeID"`
	// Styles
	Styles [4]byte `json:"styles"`
	// Lightofs is offset into lighting lump (-1 if unlit)
	Lightofs int32 `json:"lightofs"`
	// Area is total size of this side (units^2)
	Area float32 `json:"area"`
	// LightmapTextureMinsInLuxels
	LightmapTextureMinsInLuxels [2]int32 `json:"lightmapTextureMinsInLuxels"`
	// LightmapTextureSizeInLuxels is XY size of lightmap data
	LightmapTextureSizeInLuxels [2]int32 `json:"lightmapTextureSizeInLuxels"`
	// OrigFace
	OrigFace int32 `json:"origFace"`
	// NumPrims
	NumPrims uint16 `json:"numPrims"`
	// FirstPrimID
	FirstPrimID uint16 `json:"firstPrimID"`
	// SmoothingGroups
	SmoothingGroups uint32 `json:"smoothingGroups"`
}

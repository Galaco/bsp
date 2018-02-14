package lumpdata


/**
	Lump 1
 */
type Plane struct {
	Normal [3]float32 // normal vector
	Distance float32  // distance from origin
	AxisType int32	  // plane axis identifier
}

/**
	Lump 2
 */
type TexData struct {
	Reflectivity [3]float32
	NameStringTableID int32
	Width int32
	Height int32
	ViewWidth int32
	ViewHeight int32
}

/**
	Lump 4
 */
type Vis struct {
	NumClusters int32
	ByteOffset [][2]int32 // slice is incorrect. Slice length must equal numClusters (validation to be added)
}

/**
	Lump 5
 */
type Node struct{
	PlaneNum int32
	Children [2]int32
	Mins [3]int16
	Maxs [3]int16
	FirstFace uint16
	NumFaces uint16
	Area int16
	Padding int16
}

/**
	Lump 6
 */
type TexInfo struct {
	TextureVecs [2][4]float32
	LightmapVecs [2][4]float32
	Flags int32
	TexData int32
}

/**
	Lump 7
 */
type Face struct {
	Planenum uint16
	Side byte
	OnNode byte
	FirstEdge int32
	NumEdges int16
	TexInfo int16
	DispInfo int16
	SurfaceFogVolumeID int16
	Styles [4]byte
	Lightofs int32
	Area float32
	LightmapTextureMinsInLuxels [2]int32
	LightmapTextureSizeInLuxels [2]int32
	OrigFace int32
	NumPrims uint16
	FirstPrimID uint16
	SmoothingGroups uint32
}

/**
	Lump 10
 */
type Leaf struct {
	Contents int32
	Cluster int16
	Area int8 // NOTE: Actually first 9 bits of a short, but not implemented
	Flags int8 // NOTE: Actually second 7 bits of a short, but not implemented
	Mins [3]int16
	Maxs [3]int16
	FirstLeafFace uint16
	NumLeafFaces uint16
	FirstLeafBrush uint16
	NumLeafBrushes uint16
	LeafWaterDataID int16
}

/**
	Lump 14
 */
type Model struct {
	Mins [3]float32
	Maxs [3]float32
	Origin [3]float32
	HeadNode int32
	FirstFace int32
	NumFaces int32
}

/**
	Lump 18
 */
type Brush struct {
	FirstSize int32
	NumSides  int32
	Contents int32
}

/**
	Lump 19
 */
type BrushSide struct {
	PlaneNum uint16
	TexInfo int16
	DispInfo int16
	Bevel int16
}

/**
	Lump 35.
	@TODO NOTE: This really needs per-game implementations to be useful, otherwise we might as well skip reading this lump entirely
 */
type GameLumpHeader struct {
	LumpCount int32
	GameLump []GameLump // Slice length must equal lumpCount. Validation to be added

}
type GameLump struct {
	Id int32
	Flags uint16
	Version uint16
	FileOffset int32
	FileLength int32
}
type StaticPropDictLump struct {
	DictEntries int32
	Name []string // Slice length must equal dictEntries. Validation to be added
}
type StaticPropLeafLump struct {
	LeafEntries int32
	Leaf []uint16 // Slice length must equal leafEntries. Validation to be added
}
type StaticPropLump struct {
	Origin [3]float32
	Angles [3]float32
	PropType uint16
	FirstLeaf uint16
	LeafCount uint16
	Solid string //unsigned char
	Flags string //unsigned char
	Skin int32
	FadeMinDist float32
	FadeMaxDist float32
	LightingOrigin [3]float32
	ForcedFadeScale float32
	MinDXLevel uint16
	MaxDXLevel uint16
	MinCPULevel string //unsigned char
	MaxCPULevel string //unsigned char
	MinGPULevel string //unsigned char
	MaxGPULevel string //unsigned char
	//DiffuseModulation [2]float //dont know if this is correct?
	//Unknown float32
	//DisableXbox360 bool
}
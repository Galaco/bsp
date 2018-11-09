package texinfo

// Texture information for a face
// (tex, uv info, and lightmap info)
type TexInfo struct {
	TextureVecsTexelsPerWorldUnits  [2][4]float32
	LightmapVecsLuxelsPerWorldUnits [2][4]float32
	Flags                           int32
	TexData                         int32
}

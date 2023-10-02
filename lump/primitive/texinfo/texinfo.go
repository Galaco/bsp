package texinfo

// TexInfo contains texture information for a face
// (tex, uv info, and lightmap info)
type TexInfo struct {
	// TextureVecsTexelsPerWorldUnits is texture scale
	TextureVecsTexelsPerWorldUnits [2][4]float32 `json:"textureVecsTexelsPerWorldUnits"` // [s/t][xyz offset]
	// LightmapVecsLuxelsPerWorldUnits is lightmap scale
	LightmapVecsLuxelsPerWorldUnits [2][4]float32 `json:"lightmapVecsLuxelsPerWorldUnits"` // [s/t][xyz offset] - length is in units of texels/area
	// Flags
	Flags int32 `json:"flags"`
	// TexData is index into TexData lump data
	TexData int32 `json:"texData"`
}

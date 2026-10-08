package lumps

// TextureInfo contains texture information for a face
// (tex, uv info, and lightmap info)
type TextureInfo struct {
	// TextureVecsTexelsPerWorldUnits is texture scale
	TextureVecsTexelsPerWorldUnits [2][4]float32
	// LightmapVecsLuxelsPerWorldUnits is lightmap scale
	LightmapVecsLuxelsPerWorldUnits [2][4]float32
	// Flags
	Flags int32
	// TexData is index into the texture data lump.
	TexData int32
}

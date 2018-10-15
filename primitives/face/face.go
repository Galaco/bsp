package face

type Face struct {
	Planenum                    uint16
	Side                        byte
	OnNode                      byte
	FirstEdge                   int32
	NumEdges                    int16
	TexInfo                     int16
	DispInfo                    int16
	SurfaceFogVolumeID          int16
	Styles                      [4]byte
	Lightofs                    int32
	Area                        float32
	LightmapTextureMinsInLuxels [2]int32
	LightmapTextureSizeInLuxels [2]int32
	OrigFace                    int32
	NumPrims                    uint16
	FirstPrimID                 uint16
	SmoothingGroups             uint32
}

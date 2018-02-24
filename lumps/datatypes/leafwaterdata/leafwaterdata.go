package leafwaterdata


type LeafWaterData struct {
	SurfaceZ float32
	MinZ float32
	SurfaceTexInfoID int16
	AlignmentPadding int16 // Because struct is 4byte aligned
}
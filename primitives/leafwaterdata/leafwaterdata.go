package leafwaterdata

type LeafWaterData struct {
	SurfaceZ         float32
	MinZ             float32
	SurfaceTexInfoID int16
	_                int16 // Because struct is 4byte aligned
}

package leafwaterdata

// LeafWaterData
type LeafWaterData struct {
	// SurfaceZ
	SurfaceZ float32
	// MinZ
	MinZ float32
	// SurfaceTexInfoID
	SurfaceTexInfoID int16
	_                int16 // Because struct is 4byte aligned
}

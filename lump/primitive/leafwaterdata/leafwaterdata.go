package leafwaterdata

// LeafWaterData
type LeafWaterData struct {
	// SurfaceZ
	SurfaceZ float32 `json:"surfaceZ"`
	// MinZ
	MinZ float32 `json:"minZ"`
	// SurfaceTexInfoID
	SurfaceTexInfoID int16 `json:"surfaceTexInfoID"`
	_                int16 // Because struct is 4byte aligned
}

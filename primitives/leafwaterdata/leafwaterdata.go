package leafwaterdata

// LeafWaterData
type LeafWaterData struct {
	// SurfaceZ
	SurfaceZ float32
	// MinZ
	MinZ float32
	// SurfaceTexInfoID
	SurfaceTexInfoID int16
	Padding          [2]byte // Preserved bytes aligning the wire record to four bytes.
}

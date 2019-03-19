package lumps

// DispLightmapSamplePosition is Lump 34: DispLightmapSamplePosition
// NOTE: This does NOT have a mapped format yet, and is readable as []byte only
type DispLightmapSamplePosition struct {
	LumpGeneric
	data []byte
}

// Unmarshall Imports this lump from raw byte data
func (lump *DispLightmapSamplePosition) Unmarshall(raw []byte, length int32) {
	lump.data = raw
	lump.LumpInfo.SetLength(length)
}

// GetData gets internal format structure data
func (lump *DispLightmapSamplePosition) GetData() []byte {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *DispLightmapSamplePosition) Marshall() ([]byte, error) {
	return lump.data, nil
}

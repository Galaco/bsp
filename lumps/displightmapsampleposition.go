package lumps

// Lump 34: DispLightmapSamplePosition
// NOTE: This does NOT have a mapped format yet, and is readable as []byte only
type DispLightmapSamplePosition struct {
	LumpGeneric
	data []byte
}

// Import this lump from raw byte data
func (lump *DispLightmapSamplePosition) FromBytes(raw []byte, length int32) {
	lump.data = raw
	lump.LumpInfo.SetLength(length)
}

// Get internal format structure data
func (lump *DispLightmapSamplePosition) GetData() []byte {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *DispLightmapSamplePosition) ToBytes() ([]byte,error) {
	return lump.data,nil
}

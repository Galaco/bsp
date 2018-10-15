package lumps

/**
Lump 34: DispLightmapSamplePosition
*/
type DispLightmapSamplePosition struct {
	LumpGeneric
	data []byte
}

func (lump *DispLightmapSamplePosition) FromBytes(raw []byte, length int32) {
	lump.data = raw
	lump.LumpInfo.SetLength(length)
}

func (lump *DispLightmapSamplePosition) GetData() []byte {
	return lump.data
}

func (lump *DispLightmapSamplePosition) ToBytes() []byte {
	return lump.data
}

package lumps

/**
	Lump 34: DispLightmapSamplePosition
 */
type DispLightmapSamplePosition struct {
	LumpGeneric
	data []byte
}

func (lump *DispLightmapSamplePosition) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}

	lump.data = raw
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump *DispLightmapSamplePosition) GetData() interface{} {
	return lump.data
}

func (lump *DispLightmapSamplePosition) ToBytes() []byte {
	return lump.data
}

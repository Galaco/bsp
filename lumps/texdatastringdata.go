package lumps

/**
	Lump 43: TexdataStringData
 */
type TexdataStringData struct {
	LumpGeneric
	data string // MAX_MAP_TEXDATA_STRING_DATA = 256000, TEXTURE_NAME_LENGTH = 128
}

func (lump *TexdataStringData) FromBytes(raw []byte, length int32) ILump {
	lump.data = string(raw)
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump *TexdataStringData) GetData() interface{} {
	return lump.data
}

func (lump *TexdataStringData) ToBytes() []byte {
	return []byte(lump.data)
}


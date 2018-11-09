package lumps

// Lump 43: TexdataStringData
type TexdataStringData struct {
	LumpGeneric
	data string // MAX_MAP_TEXDATA_STRING_DATA = 256000, TEXTURE_NAME_LENGTH = 128
}

// Import this lump from raw byte data
func (lump *TexdataStringData) FromBytes(raw []byte, length int32) {
	lump.data = string(raw)
	lump.LumpInfo.SetLength(length)
}

// Get internal format structure data
func (lump *TexdataStringData) GetData() string {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *TexdataStringData) ToBytes() []byte {
	return []byte(lump.data)
}

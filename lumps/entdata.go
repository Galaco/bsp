package lumps

// Lump 0: Entdata
type EntData struct {
	LumpGeneric
	data string
}

// Import this lump from raw byte data
func (lump *EntData) FromBytes(raw []byte, length int32) {
	lump.data = string(raw)
	lump.LumpInfo.SetLength(length)
}

// Get internal format structure data
func (lump *EntData) GetData() string {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *EntData) ToBytes() ([]byte,error) {
	return []byte(lump.data), nil
}

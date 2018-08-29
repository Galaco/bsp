package lumps

/**
	Lump 0: Entdata
 */
type EntData struct {
	LumpGeneric
	data string
}

func (lump *EntData) FromBytes(raw []byte, length int32) {
	lump.data = string(raw)
	lump.LumpInfo.SetLength(length)
}

func (lump *EntData) GetData() string {
	return lump.data
}

func (lump *EntData) ToBytes() []byte {
	return []byte(lump.data)
}

package lumps

/**
	Lump 0: Entdata
 */
type EntData struct {
	LumpInfo
	data string
}

func (lump EntData) FromBytes(raw []byte, length int32) ILump {
	lump.data = string(raw)
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump EntData) GetData() interface{} {
	return &lump.data
}

func (lump EntData) ToBytes() []byte {
	return []byte(lump.data)
}

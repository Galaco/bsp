package lumps

/**
	Lump 28: PhysDisp
 */
type PhysDisp struct {
	LumpInfo
	data []byte
}

func (lump PhysDisp) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}

	lump.data = raw
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump PhysDisp) GetData() interface{} {
	return &lump.data
}

func (lump PhysDisp) ToBytes() []byte {
	return lump.data
}

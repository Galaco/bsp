package lumps

/**
Lump 28: PhysDisp
*/
type PhysDisp struct {
	LumpGeneric
	data []byte
}

func (lump *PhysDisp) FromBytes(raw []byte, length int32) {
	lump.data = raw
	lump.LumpInfo.SetLength(length)
}

func (lump *PhysDisp) GetData() []byte {
	return lump.data
}

func (lump *PhysDisp) ToBytes() []byte {
	return lump.data
}

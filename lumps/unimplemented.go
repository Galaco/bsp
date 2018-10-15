package lumps

/**
Lump n:
*/
type Unimplemented struct {
	LumpGeneric
	data []byte
}

func (lump *Unimplemented) FromBytes(raw []byte, length int32) {
	lump.data = raw
	lump.LumpInfo.SetLength(length)
}

func (lump *Unimplemented) GetData() []byte {
	return lump.data
}

func (lump *Unimplemented) ToBytes() []byte {
	return lump.data
}

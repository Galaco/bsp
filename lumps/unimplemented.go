package lumps


/**
	Lump n:
 */
type Unimplemented struct {
	LumpInfo
	data []byte
}

func (lump Unimplemented) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}
	lump.data = raw
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump Unimplemented) GetData() interface{} {
	return lump.data
}

func (lump Unimplemented) ToBytes() []byte {
	return lump.data
}

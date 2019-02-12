package lumps

// Lump n: Unimplemented lump type
type Unimplemented struct {
	LumpGeneric
	data []byte
}

// Import this lump from raw byte data
func (lump *Unimplemented) FromBytes(raw []byte, length int32) {
	lump.data = raw
	lump.LumpInfo.SetLength(length)
}

// Get internal format structure data
func (lump *Unimplemented) GetData() []byte {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *Unimplemented) ToBytes() ([]byte,error) {
	return lump.data,nil
}

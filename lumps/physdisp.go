package lumps

// Lump 28: PhysDisp
type PhysDisp struct {
	LumpGeneric
	data []byte
}

// Import this lump from raw byte data
func (lump *PhysDisp) Unmarshall(raw []byte, length int32) {
	lump.data = raw
	lump.LumpInfo.SetLength(length)
}

// Get internal format structure data
func (lump *PhysDisp) GetData() []byte {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *PhysDisp) Marshall() ([]byte,error) {
	return lump.data, nil
}

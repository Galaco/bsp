package lumps

// PhysDisp is Lump 28: PhysDisp
type PhysDisp struct {
	LumpGeneric
	data []byte
}

// Unmarshall Imports this lump from raw byte data
func (lump *PhysDisp) Unmarshall(raw []byte, length int32) {
	lump.data = raw
	lump.LumpInfo.SetLength(length)
}

// GetData gets internal format structure data
func (lump *PhysDisp) GetData() []byte {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *PhysDisp) Marshall() ([]byte, error) {
	return lump.data, nil
}

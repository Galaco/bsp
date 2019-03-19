package lumps

// Unimplemented is Lump n: Unimplemented lump type
type Unimplemented struct {
	LumpGeneric
	data []byte
}

// Unmarshall Imports this lump from raw byte data
func (lump *Unimplemented) Unmarshall(raw []byte, length int32) {
	lump.data = raw
	lump.LumpInfo.SetLength(length)
}

// GetData gets internal format structure data
func (lump *Unimplemented) GetData() []byte {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *Unimplemented) Marshall() ([]byte, error) {
	return lump.data, nil
}

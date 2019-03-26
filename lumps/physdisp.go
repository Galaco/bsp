package lumps

// PhysDisp is Lump 28: PhysDisp
type PhysDisp struct {
	Generic
	data []byte
}

// Unmarshall Imports this lump from raw byte data
func (lump *PhysDisp) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.data = raw
	lump.Metadata.SetLength(length)

	return err
}

// GetData gets internal format structure data
func (lump *PhysDisp) GetData() []byte {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *PhysDisp) Marshall() ([]byte, error) {
	return lump.data, nil
}

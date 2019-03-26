package lumps

// Unimplemented is Lump n: Unimplemented lump type
type Unimplemented struct {
	Generic
	data []byte
}

// Unmarshall Imports this lump from raw byte data
func (lump *Unimplemented) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.data = raw
	lump.Metadata.SetLength(length)

	return err
}

// GetData gets internal format structure data
func (lump *Unimplemented) GetData() []byte {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *Unimplemented) Marshall() ([]byte, error) {
	return lump.data, nil
}

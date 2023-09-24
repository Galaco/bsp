package lumps

// EntData is Lump 0: Entdata
type EntData struct {
	Metadata
	data string
}

// Unmarshall Imports this lump from raw byte data
func (lump *EntData) Unmarshall(raw []byte) (err error) {
	lump.data = string(raw)
	length := len(raw)
	lump.Metadata.SetLength(length)

	return err
}

// GetData gets internal format structure data
func (lump *EntData) GetData() string {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *EntData) Marshall() ([]byte, error) {
	return []byte(lump.data), nil
}

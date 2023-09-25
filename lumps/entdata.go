package lumps

// EntData is Lump 0: Entdata
type EntData struct {
	Metadata
	data string
}

// FromBytes imports this lump from raw byte data
func (lump *EntData) FromBytes(raw []byte) (err error) {
	lump.data = string(raw)
	length := len(raw)
	lump.Metadata.SetLength(length)

	return err
}

// Contents returns internal format structure data
func (lump *EntData) Contents() string {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *EntData) ToBytes() ([]byte, error) {
	return []byte(lump.data), nil
}

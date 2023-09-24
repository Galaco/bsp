package lumps

// TexDataStringData is Lump 43: TexDataStringData
type TexDataStringData struct {
	Metadata
	data string // MAX_MAP_TEXDATA_STRING_DATA = 256000, TEXTURE_NAME_LENGTH = 128
}

// Unmarshall Imports this lump from raw byte data
func (lump *TexDataStringData) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.data = string(raw)
	lump.Metadata.SetLength(length)

	return err
}

// GetData gets internal format structure data
func (lump *TexDataStringData) GetData() string {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *TexDataStringData) Marshall() ([]byte, error) {
	return []byte(lump.data), nil
}

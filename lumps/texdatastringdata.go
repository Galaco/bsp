package lumps

// TexDataStringData is Lump 43: TexDataStringData
type TexDataStringData struct {
	Metadata
	data string // MAX_MAP_TEXDATA_STRING_DATA = 256000, TEXTURE_NAME_LENGTH = 128
}

// FromBytes imports this lump from raw byte data
func (lump *TexDataStringData) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.data = string(raw)
	lump.Metadata.SetLength(length)

	return err
}

// Contents returns internal format structure data
func (lump *TexDataStringData) Contents() string {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *TexDataStringData) ToBytes() ([]byte, error) {
	return []byte(lump.data), nil
}

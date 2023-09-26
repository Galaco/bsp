package lumps

// TexDataStringData is Lump 43: TexDataStringData
type TexDataStringData struct {
	Metadata
	Data string `json:"data"` // MAX_MAP_TEXDATA_STRING_DATA = 256000, TEXTURE_NAME_LENGTH = 128
}

// FromBytes imports this lump from raw byte Data
func (lump *TexDataStringData) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Data = string(raw)
	lump.Metadata.SetLength(length)

	return err
}

// Contents returns internal format structure Data
func (lump *TexDataStringData) Contents() string {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *TexDataStringData) ToBytes() ([]byte, error) {
	return []byte(lump.Data), nil
}

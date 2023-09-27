package lump

// TexDataStringTable is Lump 44: TexDataStringTable
type TexDataStringTable struct {
	Metadata
	Data []int32 `json:"data"` // MAX_MAP_TEXINFO = 2048
}

// FromBytes imports this lump from raw byte Data
func (lump *TexDataStringTable) FromBytes(raw []byte) error {
	meta, data, err := unmarshallBasicLump[int32](raw)
	if err != nil {
		return err
	}
	lump.Metadata = meta
	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *TexDataStringTable) Contents() []int32 {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *TexDataStringTable) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

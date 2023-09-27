package lump

// VertNormalIndice is Lump 31: VertNormalIndice
type VertNormalIndice struct {
	Metadata
	Data []uint16 `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *VertNormalIndice) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[uint16](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *VertNormalIndice) Contents() []uint16 {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *VertNormalIndice) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

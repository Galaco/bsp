package lumps

// VertNormalIndice is Lump 31: VertNormalIndice
type VertNormalIndice struct {
	Metadata
	data []uint16
}

// FromBytes imports this lump from raw byte data
func (lump *VertNormalIndice) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[uint16](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *VertNormalIndice) Contents() []uint16 {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *VertNormalIndice) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

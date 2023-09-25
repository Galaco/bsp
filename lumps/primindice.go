package lumps

// PrimIndice is Lump 39: PrimIndice
type PrimIndice struct {
	Metadata
	data []uint16
}

// FromBytes imports this lump from raw byte data
func (lump *PrimIndice) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[uint16](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *PrimIndice) Contents() []uint16 {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *PrimIndice) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

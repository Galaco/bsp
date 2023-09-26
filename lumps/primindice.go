package lumps

// PrimIndice is Lump 39: PrimIndice
type PrimIndice struct {
	Metadata
	Data []uint16 `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *PrimIndice) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[uint16](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *PrimIndice) Contents() []uint16 {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *PrimIndice) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

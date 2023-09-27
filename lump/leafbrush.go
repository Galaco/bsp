package lump

// LeafBrush is Lump 17: LeafBrush
type LeafBrush struct {
	Metadata
	Data []uint16 `json:"data"` // MAX_MAP_LEAFBRUSHES = 65536
}

// FromBytes imports this lump from raw byte Data
func (lump *LeafBrush) FromBytes(raw []byte) error {
	meta, data, err := unmarshallBasicLump[uint16](raw)
	if err != nil {
		return err
	}

	lump.Metadata = meta
	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *LeafBrush) Contents() []uint16 {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *LeafBrush) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

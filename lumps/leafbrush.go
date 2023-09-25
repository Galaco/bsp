package lumps

// LeafBrush is Lump 17: LeafBrush
type LeafBrush struct {
	Metadata
	data []uint16 // MAX_MAP_LEAFBRUSHES = 65536
}

// FromBytes imports this lump from raw byte data
func (lump *LeafBrush) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[uint16](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *LeafBrush) Contents() []uint16 {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *LeafBrush) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

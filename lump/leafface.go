package lump

// LeafFace is Lump 16: LeafFace
type LeafFace struct {
	Metadata
	Data []uint16 `json:"data"` // MAX_MAP_LEAFFACES = 65536
}

// FromBytes imports this lump from raw byte Data
func (lump *LeafFace) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[uint16](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *LeafFace) Contents() []uint16 {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *LeafFace) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

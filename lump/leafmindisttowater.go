package lump

// LeafMinDistToWater is Lump 46: LeafMinDistToWater
type LeafMinDistToWater struct {
	Metadata
	Data []uint16 `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *LeafMinDistToWater) FromBytes(raw []byte) error {
	meta, data, err := unmarshallBasicLump[uint16](raw)
	if err != nil {
		return err
	}
	lump.Metadata = meta
	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *LeafMinDistToWater) Contents() []uint16 {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *LeafMinDistToWater) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

package lump

// Surfedge is Lump 13: Surfedge
type Surfedge struct {
	Metadata
	Data []int32 `json:"data"` // MAX_MAP_SURFEDGES = 512000
}

// FromBytes imports this lump from raw byte Data
func (lump *Surfedge) FromBytes(raw []byte) error {
	meta, data, err := unmarshallBasicLump[int32](raw)
	if err != nil {
		return err
	}
	lump.Metadata = meta
	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *Surfedge) Contents() []int32 {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *Surfedge) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

package lumps

// Surfedge is Lump 13: Surfedge
type Surfedge struct {
	Metadata
	data []int32 // MAX_MAP_SURFEDGES = 512000
}

// FromBytes imports this lump from raw byte data
func (lump *Surfedge) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[int32](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *Surfedge) Contents() []int32 {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Surfedge) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

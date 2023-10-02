package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/leafambientindex"
)

// LeafAmbientIndex is Lump 52: Leaf Ambient Index
type LeafAmbientIndex struct {
	Metadata
	Data []primitives.LeafAmbientIndex `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *LeafAmbientIndex) FromBytes(raw []byte) error {
	data, err := unmarshallBasicLump[primitives.LeafAmbientIndex](raw)
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *LeafAmbientIndex) Contents() []primitives.LeafAmbientIndex {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *LeafAmbientIndex) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

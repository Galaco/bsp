package lumps

import (
	primitives "github.com/galaco/bsp/primitives/leafambientindex"
)

// LeafAmbientIndex is Lump 52: Leaf Ambient Index
type LeafAmbientIndex struct {
	Metadata
	data []primitives.LeafAmbientIndex
}

// FromBytes imports this lump from raw byte data
func (lump *LeafAmbientIndex) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.LeafAmbientIndex](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *LeafAmbientIndex) Contents() []primitives.LeafAmbientIndex {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *LeafAmbientIndex) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

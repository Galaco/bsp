package lumps

import (
	primitives "github.com/galaco/bsp/primitives/leafambientlighting"
)

// LeafAmbientLighting is Lump 56: LeafAmbientLighting
type LeafAmbientLighting struct {
	Metadata
	data []primitives.LeafAmbientLighting
}

// FromBytes imports this lump from raw byte data
func (lump *LeafAmbientLighting) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.LeafAmbientLighting](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *LeafAmbientLighting) Contents() []primitives.LeafAmbientLighting {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *LeafAmbientLighting) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

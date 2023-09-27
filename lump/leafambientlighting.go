package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/leafambientlighting"
)

// LeafAmbientLighting is Lump 56: LeafAmbientLighting
type LeafAmbientLighting struct {
	Metadata
	Data []primitives.LeafAmbientLighting `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *LeafAmbientLighting) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.LeafAmbientLighting](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *LeafAmbientLighting) Contents() []primitives.LeafAmbientLighting {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *LeafAmbientLighting) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

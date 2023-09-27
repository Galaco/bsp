package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/leafambientlighting"
)

// LeafAmbientLightingHDR is Lump 55: LeafAmbientLightingHDR
type LeafAmbientLightingHDR struct {
	Metadata
	Data []primitives.LeafAmbientLighting `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *LeafAmbientLightingHDR) FromBytes(raw []byte) error {
	meta, data, err := unmarshallBasicLump[primitives.LeafAmbientLighting](raw)
	if err != nil {
		return err
	}
	lump.Metadata = meta
	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *LeafAmbientLightingHDR) Contents() []primitives.LeafAmbientLighting {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *LeafAmbientLightingHDR) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

package lumps

import (
	primitives "github.com/galaco/bsp/primitives/leafambientlighting"
)

// LeafAmbientLightingHDR is Lump 55: LeafAmbientLightingHDR
type LeafAmbientLightingHDR struct {
	Metadata
	Data []primitives.LeafAmbientLighting `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *LeafAmbientLightingHDR) FromBytes(raw []byte) error {
	meta, data, err := unmarshallBasicLump[primitives.LeafAmbientLighting](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

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

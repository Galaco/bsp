package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/leafambientindex"
)

// LeafAmbientIndexHDR is Lump 51: Leaf Ambient Index HDR
type LeafAmbientIndexHDR struct {
	Metadata
	Data []primitives.LeafAmbientIndex `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *LeafAmbientIndexHDR) FromBytes(raw []byte) error {
	data, err := unmarshallBasicLump[primitives.LeafAmbientIndex](raw)
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *LeafAmbientIndexHDR) Contents() []primitives.LeafAmbientIndex {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *LeafAmbientIndexHDR) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

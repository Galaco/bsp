package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/worldlight"
)

// WorldLight is Lump 15: Worldlight
type WorldLight struct {
	Metadata
	Data []primitives.WorldLight `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *WorldLight) FromBytes(raw []byte) error {
	data, err := unmarshallBasicLump[primitives.WorldLight](raw)
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *WorldLight) Contents() []primitives.WorldLight {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *WorldLight) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

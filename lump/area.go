package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/area"
)

// Area is Lump 20: Areas
type Area struct {
	Metadata
	Data []primitives.Area `json:"data"`
}

// FromBytes imports this lump from raw byte data.
func (lump *Area) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.Area](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *Area) Contents() []primitives.Area {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *Area) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

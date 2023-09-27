package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/brush"
)

// Brush is Lump 18: Brush
type Brush struct {
	Metadata
	Data []primitives.Brush `json:"data"`
}

// FromBytes imports this lump from raw byte data.
func (lump *Brush) FromBytes(raw []byte) error {
	meta, data, err := unmarshallBasicLump[primitives.Brush](raw)
	if err != nil {
		return err
	}
	lump.Metadata = meta
	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *Brush) Contents() []primitives.Brush {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *Brush) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

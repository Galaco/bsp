package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/common"
)

// Lighting is Lump 8: Lighting
type Lighting struct {
	Metadata
	Data []primitives.ColorRGBExponent32 `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *Lighting) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.ColorRGBExponent32](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *Lighting) Contents() []primitives.ColorRGBExponent32 {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *Lighting) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

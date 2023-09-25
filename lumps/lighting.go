package lumps

import (
	primitives "github.com/galaco/bsp/primitives/common"
)

// Lighting is Lump 8: Lighting
type Lighting struct {
	Metadata
	data []primitives.ColorRGBExponent32
}

// FromBytes imports this lump from raw byte data
func (lump *Lighting) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.ColorRGBExponent32](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *Lighting) Contents() []primitives.ColorRGBExponent32 {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Lighting) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

package lumps

import (
	primitives "github.com/galaco/bsp/primitives/primitive"
)

// Primitive is Lump 36: Primitive
type Primitive struct {
	Metadata
	data []primitives.Primitive
}

// FromBytes imports this lump from raw byte data
func (lump *Primitive) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.Primitive](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *Primitive) Contents() []primitives.Primitive {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Primitive) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

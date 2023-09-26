package lumps

import (
	primitives "github.com/galaco/bsp/primitives/primitive"
)

// Primitive is Lump 36: Primitive
type Primitive struct {
	Metadata
	Data []primitives.Primitive `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *Primitive) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.Primitive](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *Primitive) Contents() []primitives.Primitive {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *Primitive) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

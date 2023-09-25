package lumps

import (
	primitives "github.com/galaco/bsp/primitives/model"
)

// Model is Lump 14: Model
type Model struct {
	Metadata
	data []primitives.Model
}

// FromBytes imports this lump from raw byte data
func (lump *Model) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.Model](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *Model) Contents() []primitives.Model {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Model) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

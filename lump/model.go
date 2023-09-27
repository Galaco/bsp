package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/model"
)

// Model is Lump 14: Model
type Model struct {
	Metadata
	Data []primitives.Model `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *Model) FromBytes(raw []byte) error {
	meta, data, err := unmarshallBasicLump[primitives.Model](raw)
	if err != nil {
		return err
	}
	lump.Metadata = meta
	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *Model) Contents() []primitives.Model {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *Model) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

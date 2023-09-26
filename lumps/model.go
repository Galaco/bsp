package lumps

import (
	primitives "github.com/galaco/bsp/primitives/model"
)

// Model is Lump 14: Model
type Model struct {
	Metadata
	Data []primitives.Model `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *Model) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.Model](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

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

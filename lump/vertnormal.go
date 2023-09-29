package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/vertnormal"
)

// VertNormal is Lump 30: VertNormal
type VertNormal struct {
	Metadata
	Data []primitives.VertNormal `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *VertNormal) FromBytes(raw []byte) error {
	data, err := unmarshallBasicLump[primitives.VertNormal](raw)
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *VertNormal) Contents() []primitives.VertNormal {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *VertNormal) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

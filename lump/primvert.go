package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/primvert"
)

// PrimVert is Lump 37: PrimVert
type PrimVert struct {
	Metadata
	Data []primitives.PrimVert `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *PrimVert) FromBytes(raw []byte) error {
	data, err := unmarshallBasicLump[primitives.PrimVert](raw)
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *PrimVert) Contents() []primitives.PrimVert {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *PrimVert) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

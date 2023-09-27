package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/texdata"
)

// TexData is Lump 2: TexData
type TexData struct {
	Metadata
	Data []primitives.TexData `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *TexData) FromBytes(raw []byte) error {
	meta, data, err := unmarshallBasicLump[primitives.TexData](raw)
	if err != nil {
		return err
	}
	lump.Metadata = meta
	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *TexData) Contents() []primitives.TexData {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *TexData) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

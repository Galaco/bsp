package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/disptris"
)

// DispTris is Lump 48: DispTris
type DispTris struct {
	Metadata
	Data []primitives.DispTri `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *DispTris) FromBytes(raw []byte) error {
	meta, data, err := unmarshallBasicLump[primitives.DispTri](raw)
	if err != nil {
		return err
	}
	lump.Metadata = meta
	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *DispTris) Contents() []primitives.DispTri {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *DispTris) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

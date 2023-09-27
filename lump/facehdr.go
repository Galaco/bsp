package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/face"
)

// FaceHDR is Lump 58: FaceHDR
type FaceHDR struct {
	Metadata
	Data []primitives.Face `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *FaceHDR) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.Face](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *FaceHDR) Contents() []primitives.Face {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *FaceHDR) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

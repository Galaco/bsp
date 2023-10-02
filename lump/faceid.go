package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/faceid"
)

// FaceId is Lump 11: FaceIds
type FaceId struct {
	Metadata
	Data []primitives.FaceId `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *FaceId) FromBytes(raw []byte) error {
	data, err := unmarshallBasicLump[primitives.FaceId](raw)
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *FaceId) Contents() []primitives.FaceId {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *FaceId) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

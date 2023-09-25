package lumps

import (
	primitives "github.com/galaco/bsp/primitives/faceid"
)

// FaceId is Lump 11: FaceIds
type FaceId struct {
	Metadata
	data []primitives.FaceId
}

// FromBytes imports this lump from raw byte data
func (lump *FaceId) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.FaceId](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *FaceId) Contents() []primitives.FaceId {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *FaceId) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

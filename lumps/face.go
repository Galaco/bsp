package lumps

import (
	primitives "github.com/galaco/bsp/primitives/face"
)

// Face is Lump 7: Face
type Face struct {
	Metadata
	data []primitives.Face
}

// FromBytes imports this lump from raw byte data
func (lump *Face) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.Face](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *Face) Contents() []primitives.Face {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Face) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

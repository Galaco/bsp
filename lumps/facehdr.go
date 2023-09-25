package lumps

import (
	primitives "github.com/galaco/bsp/primitives/face"
)

// FaceHDR is Lump 58: FaceHDR
type FaceHDR struct {
	Metadata
	data []primitives.Face
}

// FromBytes imports this lump from raw byte data
func (lump *FaceHDR) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.Face](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *FaceHDR) Contents() []primitives.Face {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *FaceHDR) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

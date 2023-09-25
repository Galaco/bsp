package lumps

import (
	primitives "github.com/galaco/bsp/primitives/vertnormal"
)

// VertNormal is Lump 30: VertNormal
type VertNormal struct {
	Metadata
	data []primitives.VertNormal
}

// FromBytes imports this lump from raw byte data
func (lump *VertNormal) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.VertNormal](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *VertNormal) Contents() []primitives.VertNormal {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *VertNormal) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

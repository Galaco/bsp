package lumps

import (
	primitives "github.com/galaco/bsp/primitives/area"
)

// Area is Lump 20: Areas
type Area struct {
	Metadata
	data []primitives.Area
}

// FromBytes imports this lump from raw byte data
func (lump *Area) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.Area](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *Area) Contents() []primitives.Area {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Area) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

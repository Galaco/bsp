package lumps

import (
	primitives "github.com/galaco/bsp/primitives/brush"
)

// Brush is Lump 18: Brush
type Brush struct {
	Metadata
	data []primitives.Brush
}

// FromBytes imports this lump from raw byte data
func (lump *Brush) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.Brush](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *Brush) Contents() []primitives.Brush {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Brush) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

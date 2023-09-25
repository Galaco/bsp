package lumps

import (
	primitives "github.com/galaco/bsp/primitives/overlay"
)

// Overlay is Lump 45: Overlay
// Consists of an array of Overlay structs
type Overlay struct {
	Metadata
	data []primitives.Overlay
}

// FromBytes imports this lump from raw byte data
func (lump *Overlay) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.Overlay](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *Overlay) Contents() []primitives.Overlay {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Overlay) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

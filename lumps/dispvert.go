package lumps

import (
	primitives "github.com/galaco/bsp/primitives/dispvert"
)

// DispVert is Lump 33: DispVert
type DispVert struct {
	Metadata
	data []primitives.DispVert
}

// FromBytes imports this lump from raw byte data
func (lump *DispVert) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.DispVert](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *DispVert) Contents() []primitives.DispVert {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *DispVert) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

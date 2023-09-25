package lumps

import (
	primitives "github.com/galaco/bsp/primitives/primvert"
)

// PrimVert is Lump 37: PrimVert
type PrimVert struct {
	Metadata
	data []primitives.PrimVert
}

// FromBytes imports this lump from raw byte data
func (lump *PrimVert) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.PrimVert](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *PrimVert) Contents() []primitives.PrimVert {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *PrimVert) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

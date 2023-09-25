package lumps

import (
	primitives "github.com/galaco/bsp/primitives/disptris"
)

// DispTris is Lump 48: DispTris
type DispTris struct {
	Metadata
	data []primitives.DispTri
}

// FromBytes imports this lump from raw byte data
func (lump *DispTris) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.DispTri](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *DispTris) Contents() []primitives.DispTri {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *DispTris) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

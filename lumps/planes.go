package lumps

import (
	primitives "github.com/galaco/bsp/primitives/plane"
)

// Planes is Lump 1: Planes
type Planes struct {
	Metadata
	data []primitives.Plane // MAP_MAX_PLANES = 65536
}

// FromBytes imports this lump from raw byte data
func (lump *Planes) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.Plane](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *Planes) Contents() []primitives.Plane {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Planes) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

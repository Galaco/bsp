package lumps

import (
	primitives "github.com/galaco/bsp/primitives/cubemap"
)

// Cubemap is Lump 42: Cubemaps
type Cubemap struct {
	Metadata
	data []primitives.CubemapSample
}

// FromBytes imports this lump from raw byte data
func (lump *Cubemap) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.CubemapSample](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *Cubemap) Contents() []primitives.CubemapSample {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Cubemap) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

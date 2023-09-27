package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/cubemap"
)

// Cubemap is Lump 42: Cubemaps
type Cubemap struct {
	Metadata
	Data []primitives.CubemapSample `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *Cubemap) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.CubemapSample](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *Cubemap) Contents() []primitives.CubemapSample {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *Cubemap) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

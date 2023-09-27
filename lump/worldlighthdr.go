package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/worldlight"
)

// WorldLightHDR is Lump 15: Worldlight
type WorldLightHDR struct {
	Metadata
	Data []primitives.WorldLightHDR `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *WorldLightHDR) FromBytes(raw []byte) error {
	meta, data, err := unmarshallBasicLump[primitives.WorldLightHDR](raw)
	if err != nil {
		return err
	}
	lump.Metadata = meta
	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *WorldLightHDR) Contents() []primitives.WorldLightHDR {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *WorldLightHDR) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

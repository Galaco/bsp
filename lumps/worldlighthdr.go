package lumps

import (
	primitives "github.com/galaco/bsp/primitives/worldlight"
)

// WorldLightHDR is Lump 15: Worldlight
type WorldLightHDR struct {
	Metadata
	Data []primitives.WorldLightHDR `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *WorldLightHDR) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.WorldLightHDR](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

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

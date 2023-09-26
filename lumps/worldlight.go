package lumps

import (
	primitives "github.com/galaco/bsp/primitives/worldlight"
)

// WorldLight is Lump 15: Worldlight
type WorldLight struct {
	Metadata
	Data []primitives.WorldLight `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *WorldLight) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.WorldLight](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *WorldLight) Contents() []primitives.WorldLight {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *WorldLight) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

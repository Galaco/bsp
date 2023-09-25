package lumps

import (
	primitives "github.com/galaco/bsp/primitives/worldlight"
)

// WorldLight is Lump 15: Worldlight
type WorldLight struct {
	Metadata
	data []primitives.WorldLight
}

// FromBytes imports this lump from raw byte data
func (lump *WorldLight) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.WorldLight](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *WorldLight) Contents() []primitives.WorldLight {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *WorldLight) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

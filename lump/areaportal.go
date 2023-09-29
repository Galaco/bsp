package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/areaportal"
)

// AreaPortal is Lump 21: Areaportals
type AreaPortal struct {
	Metadata
	Data []primitives.AreaPortal `json:"data"`
}

// FromBytes imports this lump from raw byte data.
func (lump *AreaPortal) FromBytes(raw []byte) error {
	data, err := unmarshallBasicLump[primitives.AreaPortal](raw)
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *AreaPortal) Contents() []primitives.AreaPortal {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *AreaPortal) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

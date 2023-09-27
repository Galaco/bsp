package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/overlayfade"
)

// OverlayFade is Lump 60: Overlayfades
type OverlayFade struct {
	Metadata
	Data []primitives.OverlayFade `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *OverlayFade) FromBytes(raw []byte) error {
	meta, data, err := unmarshallBasicLump[primitives.OverlayFade](raw)
	if err != nil {
		return err
	}
	lump.Metadata = meta
	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *OverlayFade) Contents() []primitives.OverlayFade {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *OverlayFade) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

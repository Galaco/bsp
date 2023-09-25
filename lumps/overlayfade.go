package lumps

import (
	primitives "github.com/galaco/bsp/primitives/overlayfade"
)

// OverlayFade is Lump 60: Overlayfades
type OverlayFade struct {
	Metadata
	data []primitives.OverlayFade
}

// FromBytes imports this lump from raw byte data
func (lump *OverlayFade) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.OverlayFade](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *OverlayFade) Contents() []primitives.OverlayFade {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *OverlayFade) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

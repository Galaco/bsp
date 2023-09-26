package lumps

import (
	primitives "github.com/galaco/bsp/primitives/overlay"
)

// Overlay is Lump 45: Overlay
// Consists of an array of Overlay structs
type Overlay struct {
	Metadata
	Data []primitives.Overlay `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *Overlay) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.Overlay](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *Overlay) Contents() []primitives.Overlay {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *Overlay) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

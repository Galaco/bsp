package lumps

import (
	primitives "github.com/galaco/bsp/primitives/dispvert"
)

// DispVert is Lump 33: DispVert
type DispVert struct {
	Metadata
	Data []primitives.DispVert `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *DispVert) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.DispVert](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *DispVert) Contents() []primitives.DispVert {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *DispVert) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

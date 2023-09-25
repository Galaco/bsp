package lumps

import (
	primitives "github.com/galaco/bsp/primitives/texdata"
)

// TexData is Lump 2: TexData
type TexData struct {
	Metadata
	data []primitives.TexData
}

// FromBytes imports this lump from raw byte data
func (lump *TexData) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.TexData](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *TexData) Contents() []primitives.TexData {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *TexData) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

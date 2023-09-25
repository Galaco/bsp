package lumps

import (
	primitives "github.com/galaco/bsp/primitives/texinfo"
)

// TexInfo is Lump 6: TexInfo
type TexInfo struct {
	Metadata
	data []primitives.TexInfo
}

// FromBytes imports this lump from raw byte data
func (lump *TexInfo) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.TexInfo](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *TexInfo) Contents() []primitives.TexInfo {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *TexInfo) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

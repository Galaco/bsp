package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/texinfo"
)

// TexInfo is Lump 6: TexInfo
type TexInfo struct {
	Metadata
	Data []primitives.TexInfo `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *TexInfo) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.TexInfo](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *TexInfo) Contents() []primitives.TexInfo {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *TexInfo) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

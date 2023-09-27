package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/facemacrotextureinfo"
)

// FaceMacroTextureInfo is Lump 47: FaceMacroTextureInfo
type FaceMacroTextureInfo struct {
	Metadata
	Data []primitives.FaceMacroTextureInfo `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *FaceMacroTextureInfo) FromBytes(raw []byte) error {
	meta, data, err := unmarshallBasicLump[primitives.FaceMacroTextureInfo](raw)
	if err != nil {
		return err
	}
	lump.Metadata = meta
	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *FaceMacroTextureInfo) Contents() []primitives.FaceMacroTextureInfo {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *FaceMacroTextureInfo) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

package lumps

import (
	primitives "github.com/galaco/bsp/primitives/facemacrotextureinfo"
)

// FaceMacroTextureInfo is Lump 47: FaceMacroTextureInfo
type FaceMacroTextureInfo struct {
	Metadata
	Data []primitives.FaceMacroTextureInfo `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *FaceMacroTextureInfo) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.FaceMacroTextureInfo](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

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

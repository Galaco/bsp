package lumps

import (
	primitives "github.com/galaco/bsp/primitives/facemacrotextureinfo"
)

// FaceMacroTextureInfo is Lump 47: FaceMacroTextureInfo
type FaceMacroTextureInfo struct {
	Metadata
	data []primitives.FaceMacroTextureInfo
}

// FromBytes imports this lump from raw byte data
func (lump *FaceMacroTextureInfo) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.FaceMacroTextureInfo](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *FaceMacroTextureInfo) Contents() []primitives.FaceMacroTextureInfo {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *FaceMacroTextureInfo) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

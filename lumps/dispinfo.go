package lumps

import (
	primitives "github.com/galaco/bsp/primitives/dispinfo"
)

// DispInfo is Lump 26: DispInfo
type DispInfo struct {
	Metadata
	data []primitives.DispInfo
}

// FromBytes imports this lump from raw byte data
func (lump *DispInfo) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.DispInfo](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *DispInfo) Contents() []primitives.DispInfo {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *DispInfo) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

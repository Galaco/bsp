package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/dispinfo"
)

// DispInfo is Lump 26: DispInfo
type DispInfo struct {
	Metadata
	Data []primitives.DispInfo `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *DispInfo) FromBytes(raw []byte) error {
	data, err := unmarshallBasicLump[primitives.DispInfo](raw)
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *DispInfo) Contents() []primitives.DispInfo {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *DispInfo) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

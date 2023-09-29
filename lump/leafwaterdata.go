package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/leafwaterdata"
)

// LeafWaterData is Lump 36: leafwaterdata
type LeafWaterData struct {
	Metadata
	Data []primitives.LeafWaterData `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *LeafWaterData) FromBytes(raw []byte) error {
	data, err := unmarshallBasicLump[primitives.LeafWaterData](raw)
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *LeafWaterData) Contents() []primitives.LeafWaterData {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *LeafWaterData) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

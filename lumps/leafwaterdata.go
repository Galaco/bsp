package lumps

import (
	primitives "github.com/galaco/bsp/primitives/leafwaterdata"
)

// LeafWaterData is Lump 36: leafwaterdata
type LeafWaterData struct {
	Metadata
	data []primitives.LeafWaterData
}

// FromBytes imports this lump from raw byte data
func (lump *LeafWaterData) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.LeafWaterData](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *LeafWaterData) Contents() []primitives.LeafWaterData {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *LeafWaterData) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

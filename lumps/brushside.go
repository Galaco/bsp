package lumps

import (
	primitives "github.com/galaco/bsp/primitives/brushside"
)

// BrushSide is Lump 19: BrushSide
type BrushSide struct {
	Metadata
	Data []primitives.BrushSide `json:"data"` // MAX_MAP_BRUSHSIDES = 65536
}

// FromBytes imports this lump from raw byte data.
func (lump *BrushSide) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.BrushSide](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *BrushSide) Contents() []primitives.BrushSide {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *BrushSide) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

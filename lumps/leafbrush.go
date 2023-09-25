package lumps

import (
	"bytes"
	"encoding/binary"
)

// LeafBrush is Lump 17: LeafBrush
type LeafBrush struct {
	Metadata
	data []uint16 // MAX_MAP_LEAFBRUSHES = 65536
}

// FromBytes imports this lump from raw byte data
func (lump *LeafBrush) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([]uint16, length/2)
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return err
}

// Contents returns internal format structure data
func (lump *LeafBrush) Contents() []uint16 {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *LeafBrush) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

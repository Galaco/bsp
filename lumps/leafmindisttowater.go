package lumps

import (
	"bytes"
	"encoding/binary"
)

// LeafMinDistToWater is Lump 46: LeafMinDistToWater
type LeafMinDistToWater struct {
	Metadata
	data []uint16
}

// FromBytes imports this lump from raw byte data
func (lump *LeafMinDistToWater) FromBytes(raw []byte) (err error) {
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
func (lump *LeafMinDistToWater) Contents() []uint16 {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *LeafMinDistToWater) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

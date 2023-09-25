package lumps

import (
	"bytes"
	"encoding/binary"
)

// LeafFace is Lump 16: LeafFace
type LeafFace struct {
	Metadata
	data []uint16 // MAX_MAP_LEAFFACES = 65536
}

// FromBytes imports this lump from raw byte data
func (lump *LeafFace) FromBytes(raw []byte) (err error) {
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
func (lump *LeafFace) Contents() []uint16 {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *LeafFace) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

package lumps

import (
	"bytes"
	"encoding/binary"
)

// Edge is Lump 12: Edge
type Edge struct {
	Metadata
	data [][2]uint16 // MAX_MAP_EDGES = 256000
}

// FromBytes imports this lump from raw byte data
func (lump *Edge) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([][2]uint16, length/4)
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return err
}

// Contents returns internal format structure data
func (lump *Edge) Contents() [][2]uint16 {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Edge) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

package lumps

import (
	"bytes"
	"encoding/binary"
)

// Edge is Lump 12: Edge
type Edge struct {
	Metadata
	Data [][2]uint16 `json:"data"` // MAX_MAP_EDGES = 256000
}

// FromBytes imports this lump from raw byte Data
func (lump *Edge) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Data = make([][2]uint16, length/4)
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.Data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return err
}

// Contents returns internal format structure Data
func (lump *Edge) Contents() [][2]uint16 {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *Edge) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

package lumps

import (
	"bytes"
	"encoding/binary"
)

// Surfedge is Lump 13: Surfedge
type Surfedge struct {
	Metadata
	data []int32 // MAX_MAP_SURFEDGES = 512000
}

// FromBytes imports this lump from raw byte data
func (lump *Surfedge) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([]int32, length/4)
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return err
}

// Contents returns internal format structure data
func (lump *Surfedge) Contents() []int32 {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Surfedge) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

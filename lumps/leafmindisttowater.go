package lumps

import (
	"bytes"
	"encoding/binary"
	"log"
)

// Lump 46: LeafMinDistToWater
type LeafMinDistToWater struct {
	LumpGeneric
	data []uint16
}

// Import this lump from raw byte data
func (lump *LeafMinDistToWater) FromBytes(raw []byte, length int32) {
	lump.data = make([]uint16, length/int32(2))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

// Get internal format structure data
func (lump *LeafMinDistToWater) GetData() []uint16 {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *LeafMinDistToWater) ToBytes() ([]byte,error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(),err
}

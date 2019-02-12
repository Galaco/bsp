package lumps

import (
	"bytes"
	"encoding/binary"
	"log"
)

// Lump 13: Surfedge
type Surfedge struct {
	LumpGeneric
	data []int32 // MAX_MAP_SURFEDGES = 512000
}

// Import this lump from raw byte data
func (lump *Surfedge) FromBytes(raw []byte, length int32) {
	lump.data = make([]int32, length/4)
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

// Get internal format structure data
func (lump *Surfedge) GetData() []int32 {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *Surfedge) ToBytes() ([]byte,error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(),err
}

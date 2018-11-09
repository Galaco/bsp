package lumps

import (
	"bytes"
	"encoding/binary"
	"log"
)

// Lump 17: LeafBrush
type LeafBrush struct {
	LumpGeneric
	data []uint16 // MAX_MAP_LEAFBRUSHES = 65536
}

// Import this lump from raw byte data
func (lump *LeafBrush) FromBytes(raw []byte, length int32) {
	lump.data = make([]uint16, length/2)
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

// Get internal format structure data
func (lump *LeafBrush) GetData() []uint16 {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *LeafBrush) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

package lumps

import (
	"bytes"
	"encoding/binary"
	"log"
)

/**
Lump 12: Edge
*/
type Edge struct {
	LumpGeneric
	data [][2]uint16 // MAX_MAP_EDGES = 256000
}

func (lump *Edge) FromBytes(raw []byte, length int32) {
	lump.data = make([][2]uint16, length/4)
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

func (lump *Edge) GetData() [][2]uint16 {
	return lump.data
}

func (lump *Edge) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

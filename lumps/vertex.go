package lumps

import (
	"encoding/binary"
	"bytes"
	"log"
)

/**
	Lump 3: Vertex
 */

type Vertex struct {
	LumpInfo
	data [][3]float32
}
func (lump Vertex) FromBytes(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump Vertex) GetData() interface{} {
	return lump.data
}

func (lump Vertex) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

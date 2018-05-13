package lumps

import (
	"encoding/binary"
	"bytes"
	"log"
)

/**
	Lump 13: Surfedge
 */
type Surfedge struct {
	LumpInfo
	data []int32 // MAX_MAP_SURFEDGES = 512000
}

func (lump Surfedge) FromBytes(raw []byte, length int32) ILump {
	lump.data = make([]int32, length/4)
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump Surfedge) GetData() interface{} {
	return &lump.data
}

func (lump Surfedge) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

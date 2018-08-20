package lumps

import (
	"encoding/binary"
	"bytes"
	"log"
)

/**
	Lump 46: LeafMinDistToWater
 */

type LeafMinDistToWater struct {
	LumpInfo
	data []uint16
}

func (lump LeafMinDistToWater) FromBytes(raw []byte, length int32) ILump {
	lump.data = make([]uint16, length/int32(2))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump LeafMinDistToWater) GetData() interface{} {
	return lump.data
}

func (lump LeafMinDistToWater) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

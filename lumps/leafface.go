package lumps

import (
	"encoding/binary"
	"bytes"
	"log"
)

/**
	Lump 16: LeafFace
 */
type LeafFace struct {
	LumpInfo
	data []uint16 // MAX_MAP_LEAFFACES = 65536
}

func (lump LeafFace) FromBytes(raw []byte, length int32) ILump {
	lump.data = make([]uint16, length/2)
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump LeafFace) GetData() interface{} {
	return &lump.data
}

func (lump LeafFace) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

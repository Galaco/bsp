package lumps

import (
	"encoding/binary"
	"bytes"
	"log"
)

/**
	Lump 39: PrimIndice
 */
type PrimIndice struct {
	LumpGeneric
	data []uint16
}

func (lump *PrimIndice) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}

	lump.data = make([]uint16, length/int32(2))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump *PrimIndice) GetData() interface{} {
	return lump.data
}

func (lump *PrimIndice) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

package lumps

import (
	"bytes"
	"encoding/binary"
	"log"
)

/**
Lump 31: VertNormalIndice
*/
type VertNormalIndice struct {
	LumpGeneric
	data []uint16
}

func (lump *VertNormalIndice) FromBytes(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]uint16, length/int32(2))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

func (lump *VertNormalIndice) GetData() []uint16 {
	return lump.data
}

func (lump *VertNormalIndice) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

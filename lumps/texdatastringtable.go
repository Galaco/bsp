package lumps

import (
	"encoding/binary"
	"bytes"
	"log"
)

/**
	Lump 44: TexDataStringTable
 */
type TexDataStringTable struct {
	LumpGeneric
	data []int32 // MAX_MAP_TEXINFO = 2048
}

func (lump *TexDataStringTable) FromBytes(raw []byte, length int32) ILump {
	lump.data = make([]int32, length/4)
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump *TexDataStringTable) GetData() interface{} {
	return lump.data
}

func (lump *TexDataStringTable) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}
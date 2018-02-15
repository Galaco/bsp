package lumps

import (
	"encoding/binary"
	"bytes"
	"log"
)

/**
	Lump 4: Visibility
 */
type Visibility struct {
	LumpInfo
	data []byte
}

func (lump Visibility) FromBytes(raw []byte, length int32) ILump {
	lump.data = make([]byte, length)
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump Visibility) GetData() interface{} {
	return lump.data
}

func (lump Visibility) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

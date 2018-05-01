package lumps

import (
	"encoding/binary"
	"bytes"
	"log"
)

/**
	Lump 40: Pakfile
 */
type Pakfile struct {
	LumpInfo
	data []byte
}

func (lump Pakfile) FromBytes(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump Pakfile) GetData() interface{} {
	return &lump.data
}

func (lump Pakfile) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}
package lumps

import (
	"github.com/galaco/bsp/lumps/lumpdata"
	"encoding/binary"
	"bytes"
	"log"
	"unsafe"
)

/**
	Lump 2: TexData
 */
type TexData struct {
	LumpInfo
	data []lumpdata.TexData
}
func (lump TexData) FromBytes(raw []byte, length int32) ILump {
	lump.data = make([]lumpdata.TexData, length/int32(unsafe.Sizeof(lumpdata.TexData{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump TexData) GetData() interface{} {
	return lump.data
}

func (lump TexData) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

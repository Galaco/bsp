package lumps

import (
	primitives "github.com/galaco/bsp/primitives/texdata"
	"encoding/binary"
	"bytes"
	"log"
	"unsafe"
)

/**
	Lump 2: TexData
 */
type TexData struct {
	LumpGeneric
	data []primitives.TexData
}
func (lump *TexData) FromBytes(raw []byte, length int32) {
	lump.data = make([]primitives.TexData, length/int32(unsafe.Sizeof(primitives.TexData{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

func (lump *TexData) GetData()  []primitives.TexData {
	return lump.data
}

func (lump *TexData) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

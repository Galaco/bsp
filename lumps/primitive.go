package lumps

import (
	primitives "github.com/galaco/bsp/primitives/primitive"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)
/**
	Lump 36: Primitive
 */
type Primitive struct {
	LumpInfo
	data []primitives.Primitive
}

func (lump Primitive) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}

	lump.data = make([]primitives.Primitive, length/int32(unsafe.Sizeof(primitives.Primitive{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump Primitive) GetData() interface{} {
	return &lump.data
}

func (lump Primitive) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

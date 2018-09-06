package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/primitive"
	"log"
	"unsafe"
)

/**
Lump 36: Primitive
*/
type Primitive struct {
	LumpGeneric
	data []primitives.Primitive
}

func (lump *Primitive) FromBytes(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.Primitive, length/int32(unsafe.Sizeof(primitives.Primitive{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

func (lump *Primitive) GetData() []primitives.Primitive {
	return lump.data
}

func (lump *Primitive) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

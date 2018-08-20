package lumps

import (
	primitives "github.com/galaco/bsp/primitives/vertnormal"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)

/**
	Lump 30: VertNormal
 */
type VertNormal struct {
	LumpGeneric
	data []primitives.VertNormal
}

func (lump *VertNormal) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}

	lump.data = make([]primitives.VertNormal, length/int32(unsafe.Sizeof(primitives.VertNormal{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump *VertNormal) GetData() interface{} {
	return lump.data
}

func (lump *VertNormal) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

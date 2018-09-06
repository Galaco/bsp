package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/vertnormal"
	"log"
	"unsafe"
)

/**
Lump 30: VertNormal
*/
type VertNormal struct {
	LumpGeneric
	data []primitives.VertNormal
}

func (lump *VertNormal) FromBytes(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.VertNormal, length/int32(unsafe.Sizeof(primitives.VertNormal{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

func (lump *VertNormal) GetData() []primitives.VertNormal {
	return lump.data
}

func (lump *VertNormal) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

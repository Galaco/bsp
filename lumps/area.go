package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/area"
	"log"
	"unsafe"
)

/**
Lump 20: Areas
*/
type Area struct {
	LumpGeneric
	data []primitives.Area
}

func (lump *Area) FromBytes(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.Area, length/int32(unsafe.Sizeof(primitives.Area{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

func (lump *Area) GetData() []primitives.Area {
	return lump.data
}

func (lump *Area) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

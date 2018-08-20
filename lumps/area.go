package lumps

import (
	primitives "github.com/galaco/bsp/primitives/area"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)
/**
	Lump 20: Areas
 */
type Area struct {
	LumpInfo
	data []primitives.Area
}

func (lump Area) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}

	lump.data = make([]primitives.Area, length/int32(unsafe.Sizeof(primitives.Area{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump Area) GetData() interface{} {
	return lump.data
}

func (lump Area) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

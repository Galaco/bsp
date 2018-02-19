package lumps

import (
	datatypes "github.com/galaco/bsp/lumps/datatypes/area"
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
	data []datatypes.Area
}

func (lump Area) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}
	lump.data = make([]datatypes.Area, length/int32(unsafe.Sizeof(datatypes.Area{})))
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

package lumps

import (
	datatypes "github.com/galaco/bsp/lumps/datatypes/worldlight"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)
/**
	Lump 15: Worldlight
 */
type WorldLight struct {
	LumpInfo
	data []datatypes.WorldLight
}

func (lump WorldLight) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}
	lump.data = make([]datatypes.WorldLight, length/int32(unsafe.Sizeof(datatypes.WorldLight{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump WorldLight) GetData() interface{} {
	return lump.data
}

func (lump WorldLight) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

package lumps

import (
	datatypes "github.com/galaco/bsp/lumps/datatypes/vertnormal"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)
/**
	Lump 30: VertNormal
 */
type VertNormal struct {
	LumpInfo
	data []datatypes.VertNormal
}

func (lump VertNormal) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}

	lump.data = make([]datatypes.VertNormal, length/int32(unsafe.Sizeof(datatypes.VertNormal{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump VertNormal) GetData() interface{} {
	return lump.data
}

func (lump VertNormal) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}
package lumps

import (
	datatypes "github.com/galaco/bsp/lumps/datatypes/dispvert"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)
/**
	Lump 33: DispVert
 */
type DispVert struct {
	LumpInfo
	data []datatypes.DispVert
}

func (lump DispVert) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}

	lump.data = make([]datatypes.DispVert, length/int32(unsafe.Sizeof(datatypes.DispVert{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump DispVert) GetData() interface{} {
	return lump.data
}

func (lump DispVert) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

package lumps

import (
	datatypes "github.com/galaco/bsp/lumps/datatypes/disptris"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)
/**
	Lump 48: DispTris
 */
type DispTris struct {
	LumpInfo
	data []datatypes.DispTri
}

func (lump DispTris) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}

	lump.data = make([]datatypes.DispTri, length/int32(unsafe.Sizeof(datatypes.DispTri{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump DispTris) GetData() interface{} {
	return lump.data
}

func (lump DispTris) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

package lumps

import (
	datatypes "github.com/galaco/bsp/lumps/datatypes/areaportal"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)
/**
	Lump 21: Areaportals
 */
type AreaPortal struct {
	LumpInfo
	data []datatypes.AreaPortal
}

func (lump AreaPortal) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}
	lump.data = make([]datatypes.AreaPortal, length/int32(unsafe.Sizeof(datatypes.AreaPortal{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump AreaPortal) GetData() interface{} {
	return lump.data
}

func (lump AreaPortal) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

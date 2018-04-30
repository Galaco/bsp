package lumps

import (
	"encoding/binary"
	"bytes"
	"log"
	"unsafe"
	"github.com/galaco/bsp/lumps/datatypes/dispinfo"
)

/**
	Lump 26: DispInfo
 */

type DispInfo struct {
	LumpInfo
	data []dispinfo.DispInfo
}

func (lump DispInfo) FromBytes(raw []byte, length int32) ILump {
	lump.data = make([]dispinfo.DispInfo, length/int32(unsafe.Sizeof(dispinfo.DispInfo{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump DispInfo) GetData() interface{} {
	return lump.data
}

func (lump DispInfo) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

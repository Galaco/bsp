package lumps

import (
	datatypes "github.com/galaco/bsp/lumps/datatypes/texinfo"
	"encoding/binary"
	"bytes"
	"log"
	"unsafe"
)

/**
	Lump 6: TexInfo
 */

type TexInfo struct {
	LumpInfo
	data []datatypes.TexInfo
}

func (lump TexInfo) FromBytes(raw []byte, length int32) ILump {
	lump.data = make([]datatypes.TexInfo, length/int32(unsafe.Sizeof(datatypes.TexInfo{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump TexInfo) GetData() interface{} {
	return lump.data
}

func (lump TexInfo) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

package lumps

import (
	datatypes "github.com/galaco/bsp/lumps/datatypes/primvert"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)
/**
	Lump 37: PrimVert
 */
type PrimVert struct {
	LumpInfo
	data []datatypes.PrimVert
}

func (lump PrimVert) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}

	lump.data = make([]datatypes.PrimVert, length/int32(unsafe.Sizeof(datatypes.PrimVert{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump PrimVert) GetData() interface{} {
	return lump.data
}

func (lump PrimVert) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

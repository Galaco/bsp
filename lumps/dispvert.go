package lumps

import (
	primitives "github.com/galaco/bsp/primitives/dispvert"
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
	data []primitives.DispVert
}

func (lump DispVert) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}

	lump.data = make([]primitives.DispVert, length/int32(unsafe.Sizeof(primitives.DispVert{})))
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

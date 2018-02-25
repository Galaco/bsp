package lumps

import (
	datatypes "github.com/galaco/bsp/lumps/datatypes/overlay"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)

// Lump 45: Overlay
// Consists of an array of Overlay structs
type Overlay struct {
	LumpInfo
	data []datatypes.Overlay
}

func (lump Overlay) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}
	lump.data = make([]datatypes.Overlay, length/int32(unsafe.Sizeof(datatypes.Overlay{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump Overlay) GetData() interface{} {
	return lump.data
}

func (lump Overlay) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

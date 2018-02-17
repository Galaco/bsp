package lumps

import (
	datatypes "github.com/galaco/bsp/lumps/datatypes/brush"
	"encoding/binary"
	"bytes"
	"log"
	"unsafe"
)

/**
	Lump 18: Brush
 */
type Brush struct {
	LumpInfo
	data []datatypes.Brush
}

func (lump Brush) FromBytes(raw []byte, length int32) ILump {
	lump.data = make([]datatypes.Brush, length/int32(unsafe.Sizeof(datatypes.Brush{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump Brush) GetData() interface{} {
	return lump.data
}

func (lump Brush) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

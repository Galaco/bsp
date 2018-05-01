package lumps

import (
	primitives "github.com/galaco/bsp/primitives/common"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)
/**
	Lump 8: Lighting
 */
type Lighting struct {
	LumpInfo
	data []primitives.ColorRGBExponent32
}

func (lump Lighting) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}
	lump.data = make([]primitives.ColorRGBExponent32, length/int32(unsafe.Sizeof(primitives.ColorRGBExponent32{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump Lighting) GetData() interface{} {
	return &lump.data
}

func (lump Lighting) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

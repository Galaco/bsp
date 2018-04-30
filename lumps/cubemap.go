package lumps

import (
	primitives "github.com/galaco/bsp/primitives/cubemap"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)
/**
	Lump 42: Cubemaps
 */
type Cubemap struct {
	LumpInfo
	data []primitives.CubemapSample
}

func (lump Cubemap) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}
	lump.data = make([]primitives.CubemapSample, length/int32(unsafe.Sizeof(primitives.CubemapSample{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump Cubemap) GetData() interface{} {
	return lump.data
}

func (lump Cubemap) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

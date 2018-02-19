package lumps

import (
	datatypes "github.com/galaco/bsp/lumps/datatypes/cubemap"
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
	data []datatypes.CubemapSample
}

func (lump Cubemap) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}
	lump.data = make([]datatypes.CubemapSample, length/int32(unsafe.Sizeof(datatypes.CubemapSample{})))
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

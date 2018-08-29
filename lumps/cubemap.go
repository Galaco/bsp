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
	LumpGeneric
	data []primitives.CubemapSample
}

func (lump *Cubemap) FromBytes(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.CubemapSample, length/int32(unsafe.Sizeof(primitives.CubemapSample{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

func (lump *Cubemap) GetData() []primitives.CubemapSample {
	return lump.data
}

func (lump *Cubemap) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/cubemap"
	"log"
	"unsafe"
)

// Lump 42: Cubemaps
type Cubemap struct {
	LumpGeneric
	data []primitives.CubemapSample
}

// Import this lump from raw byte data
func (lump *Cubemap) FromBytes(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.CubemapSample, length/int32(unsafe.Sizeof(primitives.CubemapSample{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

// Get internal format structure data
func (lump *Cubemap) GetData() []primitives.CubemapSample {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *Cubemap) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

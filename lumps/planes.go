package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/plane"
	"log"
	"unsafe"
)

/**
Lump 1: Planes
*/
type Planes struct {
	LumpGeneric
	data []primitives.Plane // MAP_MAX_PLANES = 65536
}

func (lump *Planes) FromBytes(raw []byte, length int32) {
	lump.data = make([]primitives.Plane, length/int32(unsafe.Sizeof(primitives.Plane{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

func (lump *Planes) GetData() []primitives.Plane {
	return lump.data
}

func (lump *Planes) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

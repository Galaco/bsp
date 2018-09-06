package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/disptris"
	"log"
	"unsafe"
)

/**
Lump 48: DispTris
*/
type DispTris struct {
	LumpGeneric
	data []primitives.DispTri
}

func (lump *DispTris) FromBytes(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.DispTri, length/int32(unsafe.Sizeof(primitives.DispTri{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

func (lump *DispTris) GetData() []primitives.DispTri {
	return lump.data
}

func (lump *DispTris) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

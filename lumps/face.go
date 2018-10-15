package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/face"
	"log"
	"unsafe"
)

/**
Lump 7: Face
*/

type Face struct {
	LumpGeneric
	data []primitives.Face
}

func (lump *Face) FromBytes(raw []byte, length int32) {
	lump.data = make([]primitives.Face, length/int32(unsafe.Sizeof(primitives.Face{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

func (lump *Face) GetData() []primitives.Face {
	return lump.data
}

func (lump *Face) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/faceid"
	"log"
	"unsafe"
)

/**
Lump 11: FaceIds
*/
type FaceId struct {
	LumpGeneric
	data []primitives.FaceId
}

func (lump *FaceId) FromBytes(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.FaceId, length/int32(unsafe.Sizeof(primitives.FaceId{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

func (lump *FaceId) GetData() []primitives.FaceId {
	return lump.data
}

func (lump *FaceId) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

package lumps


import (
	primitives "github.com/galaco/bsp/primitives/faceid"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)
/**
	Lump 11: FaceIds
 */
type FaceId struct {
	LumpInfo
	data []primitives.FaceId
}

func (lump FaceId) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}
	lump.data = make([]primitives.FaceId, length/int32(unsafe.Sizeof(primitives.FaceId{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump FaceId) GetData() interface{} {
	return &lump.data
}

func (lump FaceId) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

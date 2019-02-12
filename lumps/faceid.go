package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/faceid"
	"log"
	"unsafe"
)

// Lump 11: FaceIds
type FaceId struct {
	LumpGeneric
	data []primitives.FaceId
}

// Import this lump from raw byte data
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

// Get internal format structure data
func (lump *FaceId) GetData() []primitives.FaceId {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *FaceId) ToBytes() ([]byte,error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(),err
}

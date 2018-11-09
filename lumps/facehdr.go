package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/face"
	"log"
	"unsafe"
)

// Lump 58: FaceHDR
type FaceHDR struct {
	LumpGeneric
	data []primitives.Face
}

// Import this lump from raw byte data
func (lump *FaceHDR) FromBytes(raw []byte, length int32) {
	lump.data = make([]primitives.Face, length/int32(unsafe.Sizeof(primitives.Face{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

// Get internal format structure data
func (lump *FaceHDR) GetData() []primitives.Face {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *FaceHDR) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

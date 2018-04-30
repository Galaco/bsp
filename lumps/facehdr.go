package lumps

import (
	primitives "github.com/galaco/bsp/primitives/face"
	"encoding/binary"
	"bytes"
	"log"
	"unsafe"
)

/**
	Lump 58: FaceHDR
 */

type FaceHDR struct {
	LumpInfo
	data []primitives.Face
}

func (lump FaceHDR) FromBytes(raw []byte, length int32) ILump {
	lump.data = make([]primitives.Face, length/int32(unsafe.Sizeof(primitives.Face{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump FaceHDR) GetData() interface{} {
	return lump.data
}

func (lump FaceHDR) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

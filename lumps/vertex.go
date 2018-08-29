package lumps

import (
	"encoding/binary"
	"bytes"
	"log"
	"unsafe"
	"github.com/go-gl/mathgl/mgl32"
)

/**
	Lump 3: Vertex
 */

type Vertex struct {
	LumpGeneric
	data []mgl32.Vec3
}
func (lump *Vertex) FromBytes(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]mgl32.Vec3, length/int32(unsafe.Sizeof(mgl32.Vec3{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

func (lump *Vertex) GetData() []mgl32.Vec3 {
	return lump.data
}

func (lump *Vertex) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

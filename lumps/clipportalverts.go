package lumps

import (
	"encoding/binary"
	"bytes"
	"log"
	"unsafe"
	"github.com/go-gl/mathgl/mgl32"
)

/**
	Lump 41: ClipPortalVerts
 */

type ClipPortalVerts struct {
	LumpGeneric
	data []mgl32.Vec3
}

func (lump *ClipPortalVerts) FromBytes(raw []byte, length int32) {
	lump.data = make([]mgl32.Vec3, length/int32(unsafe.Sizeof(mgl32.Vec3{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

func (lump *ClipPortalVerts) GetData() []mgl32.Vec3 {
	return lump.data
}

func (lump *ClipPortalVerts) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

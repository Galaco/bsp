package lumps

import (
	"bytes"
	"encoding/binary"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"unsafe"
)

// ClipPortalVerts is Lump 41: ClipPortalVerts
type ClipPortalVerts struct {
	LumpGeneric
	data []mgl32.Vec3
}

// Unmarshall Imports this lump from raw byte data
func (lump *ClipPortalVerts) Unmarshall(raw []byte, length int32) {
	lump.data = make([]mgl32.Vec3, length/int32(unsafe.Sizeof(mgl32.Vec3{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

// GetData gets internal format structure data
func (lump *ClipPortalVerts) GetData() []mgl32.Vec3 {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *ClipPortalVerts) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

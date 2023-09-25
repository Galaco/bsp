package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/go-gl/mathgl/mgl32"
)

// Vertex is Lump 3: Vertex
type Vertex struct {
	Metadata
	data []mgl32.Vec3
}

// FromBytes imports this lump from raw byte data
func (lump *Vertex) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]mgl32.Vec3, length/int(unsafe.Sizeof(mgl32.Vec3{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	return err
}

// Contents returns internal format structure data
func (lump *Vertex) Contents() []mgl32.Vec3 {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Vertex) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

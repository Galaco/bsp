package lumps

import (
	"bytes"
	"encoding/binary"
	"github.com/go-gl/mathgl/mgl32"
	"unsafe"
)

// Vertex is Lump 3: Vertex
type Vertex struct {
	Generic
	data []mgl32.Vec3
}

// Unmarshall Imports this lump from raw byte data
func (lump *Vertex) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]mgl32.Vec3, length/int(unsafe.Sizeof(mgl32.Vec3{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	return err
}

// GetData gets internal format structure data
func (lump *Vertex) GetData() []mgl32.Vec3 {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *Vertex) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

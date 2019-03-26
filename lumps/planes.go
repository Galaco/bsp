package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/plane"
	"unsafe"
)

// Planes is Lump 1: Planes
type Planes struct {
	Generic
	data []primitives.Plane // MAP_MAX_PLANES = 65536
}

// Unmarshall Imports this lump from raw byte data
func (lump *Planes) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([]primitives.Plane, length/int(unsafe.Sizeof(primitives.Plane{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return err
}

// GetData gets internal format structure data
func (lump *Planes) GetData() []primitives.Plane {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *Planes) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

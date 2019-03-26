package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/primitive"
	"unsafe"
)

// Primitive is Lump 36: Primitive
type Primitive struct {
	Generic
	data []primitives.Primitive
}

// Unmarshall Imports this lump from raw byte data
func (lump *Primitive) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.Primitive, length/int(unsafe.Sizeof(primitives.Primitive{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// GetData gets internal format structure data
func (lump *Primitive) GetData() []primitives.Primitive {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *Primitive) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

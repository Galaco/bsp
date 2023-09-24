package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/dispvert"
)

// DispVert is Lump 33: DispVert
type DispVert struct {
	Metadata
	data []primitives.DispVert
}

// Unmarshall Imports this lump from raw byte data
func (lump *DispVert) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.DispVert, length/int(unsafe.Sizeof(primitives.DispVert{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// GetData gets internal format structure data
func (lump *DispVert) GetData() []primitives.DispVert {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *DispVert) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

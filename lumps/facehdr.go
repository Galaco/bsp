package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/face"
	"unsafe"
)

// FaceHDR is Lump 58: FaceHDR
type FaceHDR struct {
	Generic
	data []primitives.Face
}

// Unmarshall Imports this lump from raw byte data
func (lump *FaceHDR) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([]primitives.Face, length/int(unsafe.Sizeof(primitives.Face{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return err
}

// GetData gets internal format structure data
func (lump *FaceHDR) GetData() []primitives.Face {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *FaceHDR) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

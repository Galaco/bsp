package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/faceid"
	"unsafe"
)

// FaceId is Lump 11: FaceIds
type FaceId struct {
	Generic
	data []primitives.FaceId
}

// Unmarshall Imports this lump from raw byte data
func (lump *FaceId) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.FaceId, length/int(unsafe.Sizeof(primitives.FaceId{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// GetData gets internal format structure data
func (lump *FaceId) GetData() []primitives.FaceId {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *FaceId) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

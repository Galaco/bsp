package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/face"
	"log"
	"unsafe"
)

// Face is Lump 7: Face
type Face struct {
	LumpGeneric
	data []primitives.Face
}

// Unmarshall Imports this lump from raw byte data
func (lump *Face) Unmarshall(raw []byte, length int32) {
	lump.data = make([]primitives.Face, length/int32(unsafe.Sizeof(primitives.Face{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

// GetData gets internal format structure data
func (lump *Face) GetData() []primitives.Face {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *Face) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

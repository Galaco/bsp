package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/model"
	"log"
	"unsafe"
)

// Model is Lump 14: Model
type Model struct {
	LumpGeneric
	data []primitives.Model
}

// Unmarshall Imports this lump from raw byte data
func (lump *Model) Unmarshall(raw []byte, length int32) {
	lump.data = make([]primitives.Model, length/int32(unsafe.Sizeof(primitives.Model{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

// GetData gets internal format structure data
func (lump *Model) GetData() []primitives.Model {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *Model) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

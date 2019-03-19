package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/vertnormal"
	"log"
	"unsafe"
)

// VertNormal is Lump 30: VertNormal
type VertNormal struct {
	LumpGeneric
	data []primitives.VertNormal
}

// Unmarshall Imports this lump from raw byte data
func (lump *VertNormal) Unmarshall(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.VertNormal, length/int32(unsafe.Sizeof(primitives.VertNormal{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

// GetData gets internal format structure data
func (lump *VertNormal) GetData() []primitives.VertNormal {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *VertNormal) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

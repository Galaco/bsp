package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/disptris"
	"log"
	"unsafe"
)

// DispTris is Lump 48: DispTris
type DispTris struct {
	LumpGeneric
	data []primitives.DispTri
}

// Unmarshall Imports this lump from raw byte data
func (lump *DispTris) Unmarshall(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.DispTri, length/int32(unsafe.Sizeof(primitives.DispTri{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

// GetData gets internal format structure data
func (lump *DispTris) GetData() []primitives.DispTri {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *DispTris) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

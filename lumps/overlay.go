package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/overlay"
	"log"
	"unsafe"
)

// Overlay is Lump 45: Overlay
// Consists of an array of Overlay structs
type Overlay struct {
	LumpGeneric
	data []primitives.Overlay
}

// Unmarshall Imports this lump from raw byte data
func (lump *Overlay) Unmarshall(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.Overlay, length/int32(unsafe.Sizeof(primitives.Overlay{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

// GetData gets internal format structure data
func (lump *Overlay) GetData() []primitives.Overlay {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *Overlay) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

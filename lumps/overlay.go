package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/overlay"
	"log"
	"unsafe"
)

// Lump 45: Overlay
// Consists of an array of Overlay structs
type Overlay struct {
	LumpGeneric
	data []primitives.Overlay
}

func (lump *Overlay) FromBytes(raw []byte, length int32) {
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

func (lump *Overlay) GetData() []primitives.Overlay {
	return lump.data
}

func (lump *Overlay) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

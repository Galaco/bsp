package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/area"
	"log"
	"unsafe"
)

// Lump 20: Areas
type Area struct {
	LumpGeneric
	data []primitives.Area
}

// Import this lump from raw byte data
func (lump *Area) Unmarshall(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.Area, length/int32(unsafe.Sizeof(primitives.Area{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

// Get internal format structure data
func (lump *Area) GetData() []primitives.Area {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *Area) Marshall() ([]byte,error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(),err
}

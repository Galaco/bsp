package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/brush"
	"log"
	"unsafe"
)

// Lump 18: Brush
type Brush struct {
	LumpGeneric
	data []primitives.Brush
}

// Import this lump from raw byte data
func (lump *Brush) Unmarshall(raw []byte, length int32) {
	lump.data = make([]primitives.Brush, length/int32(unsafe.Sizeof(primitives.Brush{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

// Get internal format structure data
func (lump *Brush) GetData() []primitives.Brush {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *Brush) Marshall() ([]byte,error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(),err
}

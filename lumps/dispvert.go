package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/dispvert"
	"log"
	"unsafe"
)

// Lump 33: DispVert
type DispVert struct {
	LumpGeneric
	data []primitives.DispVert
}

// Import this lump from raw byte data
func (lump *DispVert) Unmarshall(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.DispVert, length/int32(unsafe.Sizeof(primitives.DispVert{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

// Get internal format structure data
func (lump *DispVert) GetData() []primitives.DispVert {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *DispVert) Marshall() ([]byte,error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(),err
}

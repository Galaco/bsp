package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/primvert"
	"log"
	"unsafe"
)

// Lump 37: PrimVert
type PrimVert struct {
	LumpGeneric
	data []primitives.PrimVert
}

// Import this lump from raw byte data
func (lump *PrimVert) Unmarshall(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.PrimVert, length/int32(unsafe.Sizeof(primitives.PrimVert{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

// Get internal format structure data
func (lump *PrimVert) GetData() []primitives.PrimVert {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *PrimVert) Marshall() ([]byte,error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(),err
}

package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/brushside"
	"log"
	"unsafe"
)

// Lump 19: BrushSide
type BrushSide struct {
	LumpGeneric
	data []primitives.BrushSide // MAX_MAP_BRUSHSIDES = 65536
}

// Import this lump from raw byte data
func (lump *BrushSide) FromBytes(raw []byte, length int32) {
	lump.data = make([]primitives.BrushSide, length/int32(unsafe.Sizeof(primitives.BrushSide{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

// Get internal format structure data
func (lump *BrushSide) GetData() []primitives.BrushSide {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *BrushSide) ToBytes() ([]byte,error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(),err
}

package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/common"
	"log"
	"unsafe"
)

// Lump 8: Lighting
type Lighting struct {
	LumpGeneric
	data []primitives.ColorRGBExponent32
}

// Import this lump from raw byte data
func (lump *Lighting) Unmarshall(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.ColorRGBExponent32, length/int32(unsafe.Sizeof(primitives.ColorRGBExponent32{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

// Get internal format structure data
func (lump *Lighting) GetData() []primitives.ColorRGBExponent32 {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *Lighting) Marshall() ([]byte,error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(),err
}

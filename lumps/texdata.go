package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/texdata"
	"log"
	"unsafe"
)

// Lump 2: TexData
type TexData struct {
	LumpGeneric
	data []primitives.TexData
}

// Import this lump from raw byte data
func (lump *TexData) FromBytes(raw []byte, length int32) {
	lump.data = make([]primitives.TexData, length/int32(unsafe.Sizeof(primitives.TexData{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

// Get internal format structure data
func (lump *TexData) GetData() []primitives.TexData {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *TexData) ToBytes() ([]byte,error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(),err
}

package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/dispinfo"
	"log"
	"unsafe"
)

// Lump 26: DispInfo
type DispInfo struct {
	LumpGeneric
	data []primitives.DispInfo
}

// Import this lump from raw byte data
func (lump *DispInfo) Unmarshall(raw []byte, length int32) {
	lump.data = make([]primitives.DispInfo, length/int32(unsafe.Sizeof(primitives.DispInfo{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

// Get internal format structure data
func (lump *DispInfo) GetData() []primitives.DispInfo {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *DispInfo) Marshall() ([]byte,error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(),err
}

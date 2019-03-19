package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/worldlight"
	"log"
	"unsafe"
)

// Lump 15: Worldlight
type WorldLight struct {
	LumpGeneric
	data []primitives.WorldLight
}

// Import this lump from raw byte data
func (lump *WorldLight) Unmarshall(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.WorldLight, length/int32(unsafe.Sizeof(primitives.WorldLight{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

// Get internal format structure data
func (lump *WorldLight) GetData() []primitives.WorldLight {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *WorldLight) Marshall() ([]byte,error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(),err
}

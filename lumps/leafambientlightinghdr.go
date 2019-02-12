package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/leafambientlighting"
	"log"
	"unsafe"
)

// Lump 55: LeafAmbientLightingHDR
type LeafAmbientLightingHDR struct {
	LumpGeneric
	data []primitives.LeafAmbientLighting
}

// Import this lump from raw byte data
func (lump *LeafAmbientLightingHDR) FromBytes(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.LeafAmbientLighting, length/int32(unsafe.Sizeof(primitives.LeafAmbientLighting{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

// Get internal format structure data
func (lump *LeafAmbientLightingHDR) GetData() []primitives.LeafAmbientLighting {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *LeafAmbientLightingHDR) ToBytes() ([]byte,error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(),err
}

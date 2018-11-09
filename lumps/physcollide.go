package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/physcollide"
	"log"
	"unsafe"
)

// Lump 20: PhysCollide
type PhysCollide struct {
	LumpGeneric
	data []primitives.PhysCollideEntry
}

// Import this lump from raw byte data
func (lump *PhysCollide) FromBytes(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.PhysCollideEntry, length/int32(unsafe.Sizeof(primitives.PhysCollideEntry{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

// Get internal format structure data
func (lump *PhysCollide) GetData() []primitives.PhysCollideEntry {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *PhysCollide) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

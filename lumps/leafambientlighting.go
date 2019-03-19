package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/leafambientlighting"
	"log"
	"unsafe"
)

// LeafAmbientLighting is Lump 56: LeafAmbientLighting
type LeafAmbientLighting struct {
	LumpGeneric
	data []primitives.LeafAmbientLighting
}

// Unmarshall Imports this lump from raw byte data
func (lump *LeafAmbientLighting) Unmarshall(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.LeafAmbientLighting, length/int32(unsafe.Sizeof(primitives.LeafAmbientLighting{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

// GetData gets internal format structure data
func (lump *LeafAmbientLighting) GetData() []primitives.LeafAmbientLighting {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *LeafAmbientLighting) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

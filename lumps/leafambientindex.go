package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/leafambientindex"
	"log"
	"unsafe"
)

/**
Lump 52: Leaf Ambient Index
*/
type LeafAmbientIndex struct {
	LumpGeneric
	data []primitives.LeafAmbientIndex
}

func (lump *LeafAmbientIndex) FromBytes(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.LeafAmbientIndex, length/int32(unsafe.Sizeof(primitives.LeafAmbientIndex{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

func (lump *LeafAmbientIndex) GetData() []primitives.LeafAmbientIndex {
	return lump.data
}

func (lump *LeafAmbientIndex) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

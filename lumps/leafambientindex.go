package lumps

import (
	primitives "github.com/galaco/bsp/primitives/leafambientindex"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)
/**
	Lump 52: Leaf Ambient Index
 */
type LeafAmbientIndex struct {
	LumpInfo
	data []primitives.LeafAmbientIndex
}

func (lump LeafAmbientIndex) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}
	lump.data = make([]primitives.LeafAmbientIndex, length/int32(unsafe.Sizeof(primitives.LeafAmbientIndex{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump LeafAmbientIndex) GetData() interface{} {
	return lump.data
}

func (lump LeafAmbientIndex) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

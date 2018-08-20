package lumps

import (
	primitives "github.com/galaco/bsp/primitives/leafambientlighting"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)

/**
	Lump n: LeafAmbientLighting
 */
type LeafAmbientLighting struct {
	LumpGeneric
	data []primitives.LeafAmbientLighting
}

func (lump *LeafAmbientLighting) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}
	lump.data = make([]primitives.LeafAmbientLighting, length/int32(unsafe.Sizeof(primitives.LeafAmbientLighting{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump *LeafAmbientLighting) GetData() interface{} {
	return lump.data
}

func (lump *LeafAmbientLighting) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

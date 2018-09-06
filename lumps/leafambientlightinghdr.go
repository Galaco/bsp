package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/leafambientlighting"
	"log"
	"unsafe"
)

/**
Lump n: LeafAmbientLightingHDR
*/
type LeafAmbientLightingHDR struct {
	LumpGeneric
	data []primitives.LeafAmbientLighting
}

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

func (lump *LeafAmbientLightingHDR) GetData() []primitives.LeafAmbientLighting {
	return lump.data
}

func (lump *LeafAmbientLightingHDR) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

package lumps

import (
	primitives "github.com/galaco/bsp/primitives/leafambientlighting"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)

/**
	Lump n: LeafAmbientLightingHDR
 */
type LeafAmbientLightingHDR struct {
	LumpInfo
	data []primitives.LeafAmbientLighting
}

func (lump LeafAmbientLightingHDR) FromBytes(raw []byte, length int32) ILump {
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

func (lump LeafAmbientLightingHDR) GetData() interface{} {
	return lump.data
}

func (lump LeafAmbientLightingHDR) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

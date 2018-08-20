package lumps

import (
	primitives "github.com/galaco/bsp/primitives/leafambientindex"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)
/**
	Lump 51: Leaf Ambient Index HDR
 */
type LeafAmbientIndexHDR struct {
	LumpGeneric
	data []primitives.LeafAmbientIndex
}

func (lump *LeafAmbientIndexHDR) FromBytes(raw []byte, length int32) ILump {
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

func (lump *LeafAmbientIndexHDR) GetData() interface{} {
	return lump.data
}

func (lump *LeafAmbientIndexHDR) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

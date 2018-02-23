package lumps

import (
	datatypes "github.com/galaco/bsp/lumps/datatypes/leafambientindex"
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
	data []datatypes.LeafAmbientIndex
}

func (lump LeafAmbientIndex) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}
	lump.data = make([]datatypes.LeafAmbientIndex, length/int32(unsafe.Sizeof(datatypes.LeafAmbientIndex{})))
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

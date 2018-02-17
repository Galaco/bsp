package lumps

import (
	datatypes "github.com/galaco/bsp/lumps/datatypes/leaf"
	"encoding/binary"
	"bytes"
	"log"
	"unsafe"
)

/**
	Lump 10: Leaf
 */
type Leaf struct {
	LumpInfo
	data []datatypes.Leaf
}

func (lump Leaf) FromBytes(raw []byte, length int32) ILump {
	lump.data = make([]datatypes.Leaf, length/int32(unsafe.Sizeof(datatypes.Leaf{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump Leaf) GetData() interface{} {
	return lump.data
}

func (lump Leaf) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

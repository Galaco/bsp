package lumps

import (
	primitives "github.com/galaco/bsp/primitives/node"
	"encoding/binary"
	"bytes"
	"log"
	"unsafe"
)

/**
	Lump 5: Node
 */
type Node struct {
	LumpInfo
	data []primitives.Node // MAP_MAX_NODES = 65536
}

func (lump Node) FromBytes(raw []byte, length int32) ILump {
	lump.data = make([]primitives.Node, length/int32(unsafe.Sizeof(primitives.Node{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump Node) GetData() interface{} {
	return lump.data
}

func (lump Node) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

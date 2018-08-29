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
	LumpGeneric
	data []primitives.Node // MAP_MAX_NODES = 65536
}

func (lump *Node) FromBytes(raw []byte, length int32) {
	lump.data = make([]primitives.Node, length/int32(unsafe.Sizeof(primitives.Node{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

func (lump *Node) GetData() []primitives.Node {
	return lump.data
}

func (lump *Node) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

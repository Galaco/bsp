package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/node"
	"log"
	"unsafe"
)

// Node is Lump 5: Node
type Node struct {
	LumpGeneric
	data []primitives.Node // MAP_MAX_NODES = 65536
}

// Unmarshall Imports this lump from raw byte data
func (lump *Node) Unmarshall(raw []byte, length int32) {
	lump.data = make([]primitives.Node, length/int32(unsafe.Sizeof(primitives.Node{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

// GetData gets internal format structure data
func (lump *Node) GetData() []primitives.Node {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *Node) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/node"
)

// Node is Lump 5: Node
type Node struct {
	Metadata
	data []primitives.Node // MAP_MAX_NODES = 65536
}

// FromBytes imports this lump from raw byte data
func (lump *Node) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([]primitives.Node, length/int(unsafe.Sizeof(primitives.Node{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return err
}

// Contents returns internal format structure data
func (lump *Node) Contents() []primitives.Node {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Node) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

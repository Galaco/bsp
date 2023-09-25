package lumps

import (
	primitives "github.com/galaco/bsp/primitives/node"
)

// Node is Lump 5: Node
type Node struct {
	Metadata
	data []primitives.Node // MAP_MAX_NODES = 65536
}

// FromBytes imports this lump from raw byte data
func (lump *Node) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[primitives.Node](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *Node) Contents() []primitives.Node {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Node) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

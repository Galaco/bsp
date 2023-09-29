package lump

import (
	primitives "github.com/galaco/bsp/lump/primitive/node"
)

// Node is Lump 5: Node
type Node struct {
	Metadata
	Data []primitives.Node `json:"data"` // MAP_MAX_NODES = 65536
}

// FromBytes imports this lump from raw byte Data
func (lump *Node) FromBytes(raw []byte) error {
	data, err := unmarshallBasicLump[primitives.Node](raw)
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *Node) Contents() []primitives.Node {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *Node) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

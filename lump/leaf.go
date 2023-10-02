package lump

import (
	"fmt"
	primitives "github.com/galaco/bsp/lump/primitive/leaf"
)

const (
	// MaxMapLeafs is the maximum number of leafs the engine can handle
	MaxMapLeafs = 65536
)

const (
	// maxBspVersionOfV0Leaf is the last bsp revision to use the old leafv0 structure
	maxBspVersionOfV0Leaf = 19
)

// Leaf is Lump 10: Leaf
type Leaf struct {
	Metadata
	Data []primitives.Leaf `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *Leaf) FromBytes(raw []byte) error {
	data, err := unmarshallTaggedLump[primitives.Leaf](raw, fmt.Sprintf("v%d", lump.Version()))
	if err != nil {
		return err
	}

	lump.Data = data

	return nil
}

// Contents returns internal format structure Data
func (lump *Leaf) Contents() []primitives.Leaf {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *Leaf) ToBytes() ([]byte, error) {
	return marshallTaggedLump[primitives.Leaf](lump.Data, fmt.Sprintf("v%d", lump.Version()))
}

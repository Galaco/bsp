package lump

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/galaco/bsp/lump/primitive/common"
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
	// There are 2 version of leaf:
	// v0 contains a light sample
	// v1 removes the light sample, and is padded by 2 bytes
	structSize := int(unsafe.Sizeof(primitives.Leaf{}))
	if lump.Version() > maxBspVersionOfV0Leaf {
		// Correct size of v1+ leafs
		structSize -= int(unsafe.Sizeof(common.CompressedLightCube{}))
	}

	lump.Data = make([]primitives.Leaf, len(raw)/structSize)
	numLeafs := len(lump.Data)
	i := 0

	for i < numLeafs {
		leafBuf := make([]byte, structSize)
		copy(leafBuf, raw[(structSize*i):(structSize*i)+structSize])
		// Pad the raw Data to the correct size of a leaf
		if lump.Version() > maxBspVersionOfV0Leaf {
			leafBuf = append(leafBuf, make([]byte, int(unsafe.Sizeof(common.CompressedLightCube{})))...)
		}
		rawLeaf := bytes.NewBuffer(leafBuf)
		if err := binary.Read(rawLeaf, binary.LittleEndian, &lump.Data[i]); err != nil {
			return err
		}
		i++
		if i > MaxMapLeafs {
			return fmt.Errorf("leaf count overflows maximum allowed size of %d", MaxMapLeafs)
		}
	}

	return nil
}

// Contents returns internal format structure Data
func (lump *Leaf) Contents() []primitives.Leaf {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *Leaf) ToBytes() ([]byte, error) {
	var buf bytes.Buffer

	switch lump.Version() {
	case 0:
		if err := binary.Write(&buf, binary.LittleEndian, lump.Data); err != nil {
			return nil, err
		}
	default:
		structSize := int(unsafe.Sizeof(primitives.Leaf{})) - int(unsafe.Sizeof(common.CompressedLightCube{}))
		// @TODO optimize this.
		for _, l := range lump.Data {
			var leafBuf bytes.Buffer
			if err := binary.Write(&leafBuf, binary.LittleEndian, l); err != nil {
				return nil, err
			}
			if err := binary.Write(&buf, binary.LittleEndian, leafBuf.Bytes()[0:structSize]); err != nil {
				return nil, err
			}
		}
	}

	return buf.Bytes(), nil
}

package lumps

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/galaco/bsp/primitives/common"
	primitives "github.com/galaco/bsp/primitives/leaf"
	"unsafe"
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
	Generic
	data []primitives.Leaf
}

// Unmarshall Imports this lump from raw byte data
func (lump *Leaf) Unmarshall(raw []byte) (err error) {
	// There are 2 version of leaf:
	// v0 contains a light sample
	// v1 removes the light sample, and is padded by 2 bytes
	structSize := int(unsafe.Sizeof(primitives.Leaf{}))
	if lump.Version() > maxBspVersionOfV0Leaf {
		// Correct size of v1+ leafs
		structSize -= int(unsafe.Sizeof(common.CompressedLightCube{}))
	}

	length := len(raw)
	lump.data = make([]primitives.Leaf, length/structSize)
	numLeafs := len(lump.data)
	i := 0

	for i < numLeafs {
		leafBuf := make([]byte, structSize)
		copy(leafBuf, raw[(structSize * i) : (structSize*i)+structSize])
		// Pad the raw data to the correct size of a leaf
		if lump.Version() > maxBspVersionOfV0Leaf {
			leafBuf = append(leafBuf, make([]byte, int(unsafe.Sizeof(common.CompressedLightCube{})))...)
		}
		rawLeaf := bytes.NewBuffer(leafBuf)
		err = binary.Read(rawLeaf, binary.LittleEndian, &lump.data[i])
		if err != nil {
			return err
		}
		i++
		if i > MaxMapLeafs {
			return fmt.Errorf("leaf count overflows maximum allowed size of %d", MaxMapLeafs)
		}
	}
	lump.SetLength(length)

	return err
}

// GetData gets internal format structure data
func (lump *Leaf) GetData() []primitives.Leaf {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *Leaf) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	var err error

	switch lump.Version() {
	case 0:
		err = binary.Write(&buf, binary.LittleEndian, lump.data)
	default:
		structSize := int(unsafe.Sizeof(primitives.Leaf{})) - int(unsafe.Sizeof(common.CompressedLightCube{}))
		for _, l := range lump.data {
			var leafBuf bytes.Buffer
			if err = binary.Write(&leafBuf, binary.LittleEndian, l); err != nil {
				return nil, err
			}
			err = binary.Write(&buf, binary.LittleEndian, leafBuf.Bytes()[0:structSize])
		}
	}

	return buf.Bytes(), err
}

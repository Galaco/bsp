package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/leaf"
	"log"
	"unsafe"
)

const (
	// MaxMapLeafs is the maximum number of leafs the engine can handle
	MaxMapLeafs = 65536
)

// Leaf is Lump 10: Leaf
type Leaf struct {
	LumpGeneric
	data []primitives.Leaf
}

// Unmarshall Imports this lump from raw byte data
func (lump *Leaf) Unmarshall(raw []byte, length int32) {
	lump.data = make([]primitives.Leaf, length/int32(unsafe.Sizeof(primitives.Leaf{})))
	structSize := int(unsafe.Sizeof(primitives.Leaf{}))
	numLeafs := len(lump.data)
	i := 0
	for i < numLeafs {
		err := binary.Read(bytes.NewBuffer(raw[(structSize*i):(structSize*i)+structSize]), binary.LittleEndian, &lump.data[i])
		if err != nil {
			log.Fatal(err)
		}
		i++
		if i > MaxMapLeafs {
			log.Fatalf("Leaf count overflows maximum allowed size of %d\n", MaxMapLeafs)
		}
	}
	lump.LumpInfo.SetLength(length)
}

// GetData gets internal format structure data
func (lump *Leaf) GetData() []primitives.Leaf {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *Leaf) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/leaf"
	"log"
	"unsafe"
)

// Lump 10: Leaf
const MAX_MAP_LEAFS = 65536

type Leaf struct {
	LumpGeneric
	data []primitives.Leaf
}

// Import this lump from raw byte data
func (lump *Leaf) FromBytes(raw []byte, length int32) {
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
		if i > MAX_MAP_LEAFS {
			log.Fatalf("Leaf count overflows maximum allowed size of %d\n", MAX_MAP_LEAFS)
		}
	}
	lump.LumpInfo.SetLength(length)
}

// Get internal format structure data
func (lump *Leaf) GetData() []primitives.Leaf {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *Leaf) ToBytes() ([]byte,error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(),err
}

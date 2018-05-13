package lumps

import (
	primitives "github.com/galaco/bsp/primitives/leaf"
	"encoding/binary"
	"bytes"
	"log"
	"unsafe"
)

/**
	Lump 10: Leaf
 */

 const MAX_MAP_LEAFS = 65536

type Leaf struct {
	LumpInfo
	data []primitives.Leaf
}

func (lump Leaf) FromBytes(raw []byte, length int32) ILump {
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

	return lump
}

func (lump Leaf) GetData() interface{} {
	return &lump.data
}

func (lump Leaf) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

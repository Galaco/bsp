package lumps

import (
	primitives "github.com/galaco/bsp/primitives/model"
	"encoding/binary"
	"bytes"
	"log"
	"unsafe"
)

/**
	Lump 14: Model
 */
type Model struct {
	LumpGeneric
	data []primitives.Model
}

func (lump *Model) FromBytes(raw []byte, length int32) {
	lump.data = make([]primitives.Model, length/int32(unsafe.Sizeof(primitives.Model{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

func (lump *Model) GetData() []primitives.Model {
	return lump.data
}

func (lump *Model) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

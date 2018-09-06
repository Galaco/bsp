package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/brush"
	"log"
	"unsafe"
)

/**
Lump 18: Brush
*/
type Brush struct {
	LumpGeneric
	data []primitives.Brush
}

func (lump *Brush) FromBytes(raw []byte, length int32) {
	lump.data = make([]primitives.Brush, length/int32(unsafe.Sizeof(primitives.Brush{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

func (lump *Brush) GetData() []primitives.Brush {
	return lump.data
}

func (lump *Brush) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

package lumps

import (
	primitives "github.com/galaco/bsp/primitives/primvert"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)

/**
	Lump 37: PrimVert
 */
type PrimVert struct {
	LumpGeneric
	data []primitives.PrimVert
}

func (lump *PrimVert) FromBytes(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.PrimVert, length/int32(unsafe.Sizeof(primitives.PrimVert{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

func (lump *PrimVert) GetData() []primitives.PrimVert {
	return lump.data
}

func (lump *PrimVert) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

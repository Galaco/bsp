package lumps

import (
	primitives "github.com/galaco/bsp/primitives/areaportal"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)
/**
	Lump 21: Areaportals
 */
type AreaPortal struct {
	LumpGeneric
	data []primitives.AreaPortal
}

func (lump *AreaPortal) FromBytes(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.AreaPortal, length/int32(unsafe.Sizeof(primitives.AreaPortal{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

func (lump *AreaPortal) GetData() []primitives.AreaPortal {
	return lump.data
}

func (lump *AreaPortal) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

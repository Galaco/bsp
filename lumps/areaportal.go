package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/areaportal"
	"log"
	"unsafe"
)

// AreaPortal is Lump 21: Areaportals
type AreaPortal struct {
	LumpGeneric
	data []primitives.AreaPortal
}

// Unmarshall Imports this lump from raw byte data
func (lump *AreaPortal) Unmarshall(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.AreaPortal, length/int32(unsafe.Sizeof(primitives.AreaPortal{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

// GetData gets internal format structure data
func (lump *AreaPortal) GetData() []primitives.AreaPortal {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *AreaPortal) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

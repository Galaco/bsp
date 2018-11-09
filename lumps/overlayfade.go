package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/overlayfade"
	"log"
	"unsafe"
)

// Lump 60: Overlayfades
type OverlayFade struct {
	LumpGeneric
	data []primitives.OverlayFade
}

// Import this lump from raw byte data
func (lump *OverlayFade) FromBytes(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.OverlayFade, length/int32(unsafe.Sizeof(primitives.OverlayFade{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

// Get internal format structure data
func (lump *OverlayFade) GetData() []primitives.OverlayFade {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *OverlayFade) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

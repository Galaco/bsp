package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/overlayfade"
)

// OverlayFade is Lump 60: Overlayfades
type OverlayFade struct {
	Metadata
	data []primitives.OverlayFade
}

// FromBytes imports this lump from raw byte data
func (lump *OverlayFade) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.OverlayFade, length/int(unsafe.Sizeof(primitives.OverlayFade{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// Contents returns internal format structure data
func (lump *OverlayFade) Contents() []primitives.OverlayFade {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *OverlayFade) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

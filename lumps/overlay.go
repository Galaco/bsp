package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/overlay"
)

// Overlay is Lump 45: Overlay
// Consists of an array of Overlay structs
type Overlay struct {
	Metadata
	data []primitives.Overlay
}

// FromBytes imports this lump from raw byte data
func (lump *Overlay) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.Overlay, length/int(unsafe.Sizeof(primitives.Overlay{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// Contents returns internal format structure data
func (lump *Overlay) Contents() []primitives.Overlay {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Overlay) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/dispvert"
)

// DispVert is Lump 33: DispVert
type DispVert struct {
	Metadata
	data []primitives.DispVert
}

// FromBytes imports this lump from raw byte data
func (lump *DispVert) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.DispVert, length/int(unsafe.Sizeof(primitives.DispVert{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// Contents returns internal format structure data
func (lump *DispVert) Contents() []primitives.DispVert {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *DispVert) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

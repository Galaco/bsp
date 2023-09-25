package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/primitive"
)

// Primitive is Lump 36: Primitive
type Primitive struct {
	Metadata
	data []primitives.Primitive
}

// FromBytes imports this lump from raw byte data
func (lump *Primitive) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.Primitive, length/int(unsafe.Sizeof(primitives.Primitive{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// Contents returns internal format structure data
func (lump *Primitive) Contents() []primitives.Primitive {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Primitive) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/primvert"
)

// PrimVert is Lump 37: PrimVert
type PrimVert struct {
	Metadata
	data []primitives.PrimVert
}

// FromBytes imports this lump from raw byte data
func (lump *PrimVert) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.PrimVert, length/int(unsafe.Sizeof(primitives.PrimVert{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// Contents returns internal format structure data
func (lump *PrimVert) Contents() []primitives.PrimVert {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *PrimVert) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/vertnormal"
)

// VertNormal is Lump 30: VertNormal
type VertNormal struct {
	Metadata
	data []primitives.VertNormal
}

// FromBytes imports this lump from raw byte data
func (lump *VertNormal) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.VertNormal, length/int(unsafe.Sizeof(primitives.VertNormal{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	return err
}

// Contents returns internal format structure data
func (lump *VertNormal) Contents() []primitives.VertNormal {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *VertNormal) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

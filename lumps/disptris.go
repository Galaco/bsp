package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/disptris"
)

// DispTris is Lump 48: DispTris
type DispTris struct {
	Metadata
	data []primitives.DispTri
}

// FromBytes imports this lump from raw byte data
func (lump *DispTris) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.DispTri, length/int(unsafe.Sizeof(primitives.DispTri{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// Contents returns internal format structure data
func (lump *DispTris) Contents() []primitives.DispTri {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *DispTris) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

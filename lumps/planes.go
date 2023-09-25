package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/plane"
)

// Planes is Lump 1: Planes
type Planes struct {
	Metadata
	data []primitives.Plane // MAP_MAX_PLANES = 65536
}

// FromBytes imports this lump from raw byte data
func (lump *Planes) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([]primitives.Plane, length/int(unsafe.Sizeof(primitives.Plane{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return err
}

// Contents returns internal format structure data
func (lump *Planes) Contents() []primitives.Plane {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Planes) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

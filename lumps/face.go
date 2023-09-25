package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/face"
)

// Face is Lump 7: Face
type Face struct {
	Metadata
	data []primitives.Face
}

// FromBytes imports this lump from raw byte data
func (lump *Face) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([]primitives.Face, length/int(unsafe.Sizeof(primitives.Face{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return err
}

// Contents returns internal format structure data
func (lump *Face) Contents() []primitives.Face {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Face) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

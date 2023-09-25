package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/faceid"
)

// FaceId is Lump 11: FaceIds
type FaceId struct {
	Metadata
	data []primitives.FaceId
}

// FromBytes imports this lump from raw byte data
func (lump *FaceId) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.FaceId, length/int(unsafe.Sizeof(primitives.FaceId{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// Contents returns internal format structure data
func (lump *FaceId) Contents() []primitives.FaceId {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *FaceId) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

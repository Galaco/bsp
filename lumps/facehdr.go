package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/face"
)

// FaceHDR is Lump 58: FaceHDR
type FaceHDR struct {
	Metadata
	data []primitives.Face
}

// FromBytes imports this lump from raw byte data
func (lump *FaceHDR) FromBytes(raw []byte) (err error) {
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
func (lump *FaceHDR) Contents() []primitives.Face {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *FaceHDR) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

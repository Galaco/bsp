package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/model"
)

// Model is Lump 14: Model
type Model struct {
	Metadata
	data []primitives.Model
}

// FromBytes imports this lump from raw byte data
func (lump *Model) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([]primitives.Model, length/int(unsafe.Sizeof(primitives.Model{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return err
}

// Contents returns internal format structure data
func (lump *Model) Contents() []primitives.Model {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Model) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

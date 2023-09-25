package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/texdata"
)

// TexData is Lump 2: TexData
type TexData struct {
	Metadata
	data []primitives.TexData
}

// FromBytes imports this lump from raw byte data
func (lump *TexData) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([]primitives.TexData, length/int(unsafe.Sizeof(primitives.TexData{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return err
}

// Contents returns internal format structure data
func (lump *TexData) Contents() []primitives.TexData {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *TexData) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

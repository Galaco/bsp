package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/brush"
)

// Brush is Lump 18: Brush
type Brush struct {
	Metadata
	data []primitives.Brush
}

// FromBytes imports this lump from raw byte data
func (lump *Brush) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([]primitives.Brush, length/int(unsafe.Sizeof(primitives.Brush{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return nil
}

// Contents returns internal format structure data
func (lump *Brush) Contents() []primitives.Brush {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Brush) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

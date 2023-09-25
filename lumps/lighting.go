package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/common"
)

// Lighting is Lump 8: Lighting
type Lighting struct {
	Metadata
	data []primitives.ColorRGBExponent32
}

// FromBytes imports this lump from raw byte data
func (lump *Lighting) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.ColorRGBExponent32, length/int(unsafe.Sizeof(primitives.ColorRGBExponent32{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// Contents returns internal format structure data
func (lump *Lighting) Contents() []primitives.ColorRGBExponent32 {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Lighting) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

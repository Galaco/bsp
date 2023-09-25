package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/worldlight"
)

// WorldLight is Lump 15: Worldlight
type WorldLight struct {
	Metadata
	data []primitives.WorldLight
}

// FromBytes imports this lump from raw byte data
func (lump *WorldLight) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.WorldLight, length/int(unsafe.Sizeof(primitives.WorldLight{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	return err
}

// Contents returns internal format structure data
func (lump *WorldLight) Contents() []primitives.WorldLight {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *WorldLight) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

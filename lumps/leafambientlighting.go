package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/leafambientlighting"
)

// LeafAmbientLighting is Lump 56: LeafAmbientLighting
type LeafAmbientLighting struct {
	Metadata
	data []primitives.LeafAmbientLighting
}

// FromBytes imports this lump from raw byte data
func (lump *LeafAmbientLighting) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.LeafAmbientLighting, length/int(unsafe.Sizeof(primitives.LeafAmbientLighting{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return err
}

// Contents returns internal format structure data
func (lump *LeafAmbientLighting) Contents() []primitives.LeafAmbientLighting {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *LeafAmbientLighting) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

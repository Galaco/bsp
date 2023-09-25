package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/leafambientlighting"
)

// LeafAmbientLightingHDR is Lump 55: LeafAmbientLightingHDR
type LeafAmbientLightingHDR struct {
	Metadata
	data []primitives.LeafAmbientLighting
}

// FromBytes imports this lump from raw byte data
func (lump *LeafAmbientLightingHDR) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.LeafAmbientLighting, length/int(unsafe.Sizeof(primitives.LeafAmbientLighting{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// Contents returns internal format structure data
func (lump *LeafAmbientLightingHDR) Contents() []primitives.LeafAmbientLighting {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *LeafAmbientLightingHDR) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

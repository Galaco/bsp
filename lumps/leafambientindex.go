package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/leafambientindex"
)

// LeafAmbientIndex is Lump 52: Leaf Ambient Index
type LeafAmbientIndex struct {
	Metadata
	data []primitives.LeafAmbientIndex
}

// FromBytes imports this lump from raw byte data
func (lump *LeafAmbientIndex) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.LeafAmbientIndex, length/int(unsafe.Sizeof(primitives.LeafAmbientIndex{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// Contents returns internal format structure data
func (lump *LeafAmbientIndex) Contents() []primitives.LeafAmbientIndex {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *LeafAmbientIndex) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/leafambientindex"
)

// LeafAmbientIndexHDR is Lump 51: Leaf Ambient Index HDR
type LeafAmbientIndexHDR struct {
	Metadata
	data []primitives.LeafAmbientIndex
}

// FromBytes imports this lump from raw byte data
func (lump *LeafAmbientIndexHDR) FromBytes(raw []byte) (err error) {
	length := len(raw)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.LeafAmbientIndex, length/int(unsafe.Sizeof(primitives.LeafAmbientIndex{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return err
}

// Contents returns internal format structure data
func (lump *LeafAmbientIndexHDR) Contents() []primitives.LeafAmbientIndex {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *LeafAmbientIndexHDR) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

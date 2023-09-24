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

// Unmarshall Imports this lump from raw byte data
func (lump *LeafAmbientIndex) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.LeafAmbientIndex, length/int(unsafe.Sizeof(primitives.LeafAmbientIndex{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// GetData gets internal format structure data
func (lump *LeafAmbientIndex) GetData() []primitives.LeafAmbientIndex {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *LeafAmbientIndex) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

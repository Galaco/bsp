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

// Unmarshall Imports this lump from raw byte data
func (lump *LeafAmbientIndexHDR) Unmarshall(raw []byte) (err error) {
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

// GetData gets internal format structure data
func (lump *LeafAmbientIndexHDR) GetData() []primitives.LeafAmbientIndex {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *LeafAmbientIndexHDR) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

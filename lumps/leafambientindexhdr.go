package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/leafambientindex"
	"log"
	"unsafe"
)

// LeafAmbientIndexHDR is Lump 51: Leaf Ambient Index HDR
type LeafAmbientIndexHDR struct {
	LumpGeneric
	data []primitives.LeafAmbientIndex
}

// Unmarshall Imports this lump from raw byte data
func (lump *LeafAmbientIndexHDR) Unmarshall(raw []byte, length int32) {
	if length == 0 {
		return
	}
	lump.data = make([]primitives.LeafAmbientIndex, length/int32(unsafe.Sizeof(primitives.LeafAmbientIndex{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
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

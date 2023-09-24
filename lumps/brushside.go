package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/brushside"
)

// BrushSide is Lump 19: BrushSide
type BrushSide struct {
	Metadata
	data []primitives.BrushSide // MAX_MAP_BRUSHSIDES = 65536
}

// Unmarshall Imports this lump from raw byte data
func (lump *BrushSide) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([]primitives.BrushSide, length/int(unsafe.Sizeof(primitives.BrushSide{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return nil
}

// GetData gets internal format structure data
func (lump *BrushSide) GetData() []primitives.BrushSide {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *BrushSide) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

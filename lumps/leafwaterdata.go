package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/leafwaterdata"
)

// LeafWaterData is Lump 36: leafwaterdata
type LeafWaterData struct {
	Metadata
	data []primitives.LeafWaterData
}

// Unmarshall Imports this lump from raw byte data
func (lump *LeafWaterData) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.LeafWaterData, length/int(unsafe.Sizeof(primitives.LeafWaterData{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// GetData gets internal format structure data
func (lump *LeafWaterData) GetData() []primitives.LeafWaterData {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *LeafWaterData) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

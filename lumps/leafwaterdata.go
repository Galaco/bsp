package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/leafwaterdata"
	"log"
	"unsafe"
)

// LeafWaterData is Lump 36: leafwaterdata
type LeafWaterData struct {
	LumpGeneric
	data []primitives.LeafWaterData
}

// Unmarshall Imports this lump from raw byte data
func (lump *LeafWaterData) Unmarshall(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.LeafWaterData, length/int32(unsafe.Sizeof(primitives.LeafWaterData{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

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

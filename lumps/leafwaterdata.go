package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/leafwaterdata"
	"log"
	"unsafe"
)

// Lump 36: leafwaterdata
type LeafWaterData struct {
	LumpGeneric
	data []primitives.LeafWaterData
}

// Import this lump from raw byte data
func (lump *LeafWaterData) FromBytes(raw []byte, length int32) {
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

// Get internal format structure data
func (lump *LeafWaterData) GetData() []primitives.LeafWaterData {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *LeafWaterData) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

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

// FromBytes imports this lump from raw byte data
func (lump *LeafWaterData) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.LeafWaterData, length/int(unsafe.Sizeof(primitives.LeafWaterData{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// Contents returns internal format structure data
func (lump *LeafWaterData) Contents() []primitives.LeafWaterData {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *LeafWaterData) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

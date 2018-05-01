package lumps

import (
	primitives "github.com/galaco/bsp/primitives/leafwaterdata"
	"encoding/binary"
	"bytes"
	"log"
	"unsafe"
)

/**
	Lump 36: leafwaterdata
 */
type LeafWaterData struct {
	LumpInfo
	data []primitives.LeafWaterData
}

func (lump LeafWaterData) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}
	lump.data = make([]primitives.LeafWaterData, length/int32(unsafe.Sizeof(primitives.LeafWaterData{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump LeafWaterData) GetData() interface{} {
	return &lump.data
}

func (lump LeafWaterData) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

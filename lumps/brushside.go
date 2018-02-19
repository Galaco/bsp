package lumps

import (
	datatypes "github.com/galaco/bsp/lumps/datatypes/brushside"
	"encoding/binary"
	"bytes"
	"log"
	"unsafe"
)

/**
	Lump 19: BrushSide
 */
type BrushSide struct {
	LumpInfo
	data []datatypes.BrushSide // MAX_MAP_BRUSHSIDES = 65536
}

func (lump BrushSide) FromBytes(raw []byte, length int32) ILump {
	lump.data = make([]datatypes.BrushSide, length/int32(unsafe.Sizeof(datatypes.BrushSide{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump BrushSide) GetData() interface{} {
	return lump.data
}

func (lump BrushSide) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

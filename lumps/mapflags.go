package lumps

import (
	primitives "github.com/galaco/bsp/primitives/mapflags"
	"encoding/binary"
	"bytes"
	"log"
)
/**
	Lump 59: MapFlags
 */
type MapFlags struct {
	LumpInfo
	data primitives.MapFlags
}

func (lump MapFlags) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}

	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump MapFlags) GetData() interface{} {
	return lump.data
}

func (lump MapFlags) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

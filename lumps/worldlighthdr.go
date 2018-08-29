package lumps

import (
	primitives "github.com/galaco/bsp/primitives/worldlight"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)

/**
	Lump 15: Worldlight
 */
type WorldLightHDR struct {
	LumpGeneric
	data []primitives.WorldLight
}

func (lump *WorldLightHDR) FromBytes(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.WorldLight, length/int32(unsafe.Sizeof(primitives.WorldLight{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

func (lump *WorldLightHDR) GetData() []primitives.WorldLight {
	return lump.data
}

func (lump *WorldLightHDR) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

package lumps

import (
	"encoding/binary"
	"bytes"
	"log"
	"unsafe"
	"github.com/galaco/bsp/primitives/common"
)

/**
	Lump 41: ClipPortalVerts
 */

type ClipPortalVerts struct {
	LumpInfo
	data []common.Vector
}

func (lump ClipPortalVerts) FromBytes(raw []byte, length int32) ILump {
	lump.data = make([]common.Vector, length/int32(unsafe.Sizeof(common.Vector{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump ClipPortalVerts) GetData() interface{} {
	return lump.data
}

func (lump ClipPortalVerts) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

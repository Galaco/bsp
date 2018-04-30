package lumps

import (
	primitives "github.com/galaco/bsp/primitives/facemacrotextureinfo"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)
/**
	Lump 47: FaceMacroTextureInfo
 */
type FaceMacroTextureInfo struct {
	LumpInfo
	data []primitives.FaceMacroTextureInfo
}

func (lump FaceMacroTextureInfo) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}

	lump.data = make([]primitives.FaceMacroTextureInfo, length/int32(unsafe.Sizeof(primitives.FaceMacroTextureInfo{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump FaceMacroTextureInfo) GetData() interface{} {
	return lump.data
}

func (lump FaceMacroTextureInfo) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

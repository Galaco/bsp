package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/facemacrotextureinfo"
	"log"
	"unsafe"
)

/**
Lump 47: FaceMacroTextureInfo
*/
type FaceMacroTextureInfo struct {
	LumpGeneric
	data []primitives.FaceMacroTextureInfo
}

func (lump *FaceMacroTextureInfo) FromBytes(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.FaceMacroTextureInfo, length/int32(unsafe.Sizeof(primitives.FaceMacroTextureInfo{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
}

func (lump *FaceMacroTextureInfo) GetData() []primitives.FaceMacroTextureInfo {
	return lump.data
}

func (lump *FaceMacroTextureInfo) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

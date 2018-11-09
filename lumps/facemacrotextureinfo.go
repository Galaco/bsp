package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/facemacrotextureinfo"
	"log"
	"unsafe"
)

// Lump 47: FaceMacroTextureInfo
type FaceMacroTextureInfo struct {
	LumpGeneric
	data []primitives.FaceMacroTextureInfo
}

// Import this lump from raw byte data
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

// Get internal format structure data
func (lump *FaceMacroTextureInfo) GetData() []primitives.FaceMacroTextureInfo {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *FaceMacroTextureInfo) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

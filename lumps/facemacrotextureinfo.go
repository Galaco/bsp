package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/facemacrotextureinfo"
)

// FaceMacroTextureInfo is Lump 47: FaceMacroTextureInfo
type FaceMacroTextureInfo struct {
	Metadata
	data []primitives.FaceMacroTextureInfo
}

// FromBytes imports this lump from raw byte data
func (lump *FaceMacroTextureInfo) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.FaceMacroTextureInfo, length/int(unsafe.Sizeof(primitives.FaceMacroTextureInfo{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// Contents returns internal format structure data
func (lump *FaceMacroTextureInfo) Contents() []primitives.FaceMacroTextureInfo {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *FaceMacroTextureInfo) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

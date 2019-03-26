package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/facemacrotextureinfo"
	"unsafe"
)

// FaceMacroTextureInfo is Lump 47: FaceMacroTextureInfo
type FaceMacroTextureInfo struct {
	Generic
	data []primitives.FaceMacroTextureInfo
}

// Unmarshall Imports this lump from raw byte data
func (lump *FaceMacroTextureInfo) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.FaceMacroTextureInfo, length/int(unsafe.Sizeof(primitives.FaceMacroTextureInfo{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// GetData gets internal format structure data
func (lump *FaceMacroTextureInfo) GetData() []primitives.FaceMacroTextureInfo {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *FaceMacroTextureInfo) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/texinfo"
)

// TexInfo is Lump 6: TexInfo
type TexInfo struct {
	Metadata
	data []primitives.TexInfo
}

// FromBytes imports this lump from raw byte data
func (lump *TexInfo) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([]primitives.TexInfo, length/int(unsafe.Sizeof(primitives.TexInfo{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return err
}

// Contents returns internal format structure data
func (lump *TexInfo) Contents() []primitives.TexInfo {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *TexInfo) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

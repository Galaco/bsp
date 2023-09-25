package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/dispinfo"
)

// DispInfo is Lump 26: DispInfo
type DispInfo struct {
	Metadata
	data []primitives.DispInfo
}

// FromBytes imports this lump from raw byte data
func (lump *DispInfo) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([]primitives.DispInfo, length/int(unsafe.Sizeof(primitives.DispInfo{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return nil
}

// Contents returns internal format structure data
func (lump *DispInfo) Contents() []primitives.DispInfo {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *DispInfo) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

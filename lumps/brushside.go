package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/brushside"
)

// BrushSide is Lump 19: BrushSide
type BrushSide struct {
	Metadata
	data []primitives.BrushSide // MAX_MAP_BRUSHSIDES = 65536
}

// FromBytes imports this lump from raw byte data
func (lump *BrushSide) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([]primitives.BrushSide, length/int(unsafe.Sizeof(primitives.BrushSide{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return nil
}

// Contents returns internal format structure data
func (lump *BrushSide) Contents() []primitives.BrushSide {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *BrushSide) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

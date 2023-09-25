package lumps

import (
	"bytes"
	"encoding/binary"
)

// TexDataStringTable is Lump 44: TexDataStringTable
type TexDataStringTable struct {
	Metadata
	data []int32 // MAX_MAP_TEXINFO = 2048
}

// FromBytes imports this lump from raw byte data
func (lump *TexDataStringTable) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([]int32, length/4)
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return nil
}

// Contents returns internal format structure data
func (lump *TexDataStringTable) Contents() []int32 {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *TexDataStringTable) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

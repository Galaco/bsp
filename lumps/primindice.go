package lumps

import (
	"bytes"
	"encoding/binary"
)

// PrimIndice is Lump 39: PrimIndice
type PrimIndice struct {
	Metadata
	data []uint16
}

// FromBytes imports this lump from raw byte data
func (lump *PrimIndice) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]uint16, length/2)
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// Contents returns internal format structure data
func (lump *PrimIndice) Contents() []uint16 {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *PrimIndice) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

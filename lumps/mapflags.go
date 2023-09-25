package lumps

import (
	"bytes"
	"encoding/binary"

	primitives "github.com/galaco/bsp/primitives/mapflags"
)

// MapFlags is Lump 59: MapFlags
type MapFlags struct {
	Metadata
	data primitives.MapFlags
}

// FromBytes imports this lump from raw byte data
func (lump *MapFlags) FromBytes(raw []byte) (err error) {
	length := len(raw)
	if length == 0 {
		return
	}

	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return err
}

// Contents returns internal format structure data
func (lump *MapFlags) Contents() *primitives.MapFlags {
	return &lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *MapFlags) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}

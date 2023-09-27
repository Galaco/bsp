package lump

import (
	"bytes"
	"encoding/binary"

	primitives "github.com/galaco/bsp/lump/primitive/mapflags"
)

// MapFlags is Lump 59: MapFlags
type MapFlags struct {
	Metadata
	Data primitives.MapFlags `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *MapFlags) FromBytes(raw []byte) error {
	if len(raw) == 0 {
		return nil
	}

	return binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.Data)
}

// Contents returns internal format structure Data
func (lump *MapFlags) Contents() *primitives.MapFlags {
	return &lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *MapFlags) ToBytes() ([]byte, error) {
	a, err := marshallBasicLump(lump.Data)
	return a, err
}

package lumps

import (
	"bytes"
	"encoding/binary"

	primitives "github.com/galaco/bsp/primitives/mapflags"
)

// MapFlags is Lump 59: MapFlags
type MapFlags struct {
	Metadata
	Data primitives.MapFlags `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *MapFlags) FromBytes(raw []byte) error {
	length := len(raw)
	if length == 0 {
		return nil
	}

	if err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.Data); err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return nil
}

// Contents returns internal format structure Data
func (lump *MapFlags) Contents() *primitives.MapFlags {
	return &lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *MapFlags) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}

package lump

import (
	"bytes"
	"encoding/binary"

	primitives "github.com/galaco/bsp/lump/primitive/visibility"
)

// Visibility is Lump 4: Visibility
type Visibility struct {
	Metadata
	Data primitives.Vis `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *Visibility) FromBytes(raw []byte) (err error) {
	length := len(raw)
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.Data.NumClusters)
	if err != nil {
		return err
	}
	lump.Data.ByteOffset = make([][2]int32, lump.Data.NumClusters)
	err = binary.Read(bytes.NewBuffer(raw[4:]), binary.LittleEndian, &lump.Data.ByteOffset)
	if err != nil {
		return err
	}
	lump.Data.BitVectors = make([]byte, length)
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.Data.BitVectors)
	if err != nil {
		return err
	}

	lump.Metadata.SetLength(length)

	return nil
}

// Contents returns internal format structure Data
func (lump *Visibility) Contents() *primitives.Vis {
	return &lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *Visibility) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data.BitVectors)
}

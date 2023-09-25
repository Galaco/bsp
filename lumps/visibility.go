package lumps

import (
	"bytes"
	"encoding/binary"

	primitives "github.com/galaco/bsp/primitives/visibility"
)

// Visibility is Lump 4: Visibility
type Visibility struct {
	Metadata
	data primitives.Vis
}

// FromBytes imports this lump from raw byte data
func (lump *Visibility) FromBytes(raw []byte) (err error) {
	length := len(raw)
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data.NumClusters)
	if err != nil {
		return err
	}
	lump.data.ByteOffset = make([][2]int32, lump.data.NumClusters)
	err = binary.Read(bytes.NewBuffer(raw[4:]), binary.LittleEndian, &lump.data.ByteOffset)
	if err != nil {
		return err
	}
	lump.data.BitVectors = make([]byte, length)
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data.BitVectors)
	if err != nil {
		return err
	}

	lump.Metadata.SetLength(length)

	return nil
}

// Contents returns internal format structure data
func (lump *Visibility) Contents() *primitives.Vis {
	return &lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Visibility) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data.BitVectors)
}

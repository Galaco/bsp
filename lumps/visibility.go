package lumps

import (
	"encoding/binary"
	"bytes"
	"log"
	primitives "github.com/galaco/bsp/primitives/visibility"
)

// Lump 4: Visibility
type Visibility struct {
	LumpGeneric
	data primitives.Vis
}

// FromBytes
// Populate receiver lump from byte slice
func (lump *Visibility) FromBytes(raw []byte, length int32) {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data.NumClusters)
	if err != nil {
		log.Fatal(err)
	}
	lump.data.ByteOffset = make([][2]int32, lump.data.NumClusters)
	err = binary.Read(bytes.NewBuffer(raw[4:]), binary.LittleEndian, &lump.data.ByteOffset)
	if err != nil {
		log.Fatal(err)
	}
	offset := 4 + (8 * lump.data.NumClusters)
	lump.data.BitVectors = make([]byte, length - offset)
	err = binary.Read(bytes.NewBuffer(raw[offset:]), binary.LittleEndian, &lump.data.BitVectors)
	if err != nil {
		log.Fatal(err)
	}

	lump.LumpInfo.SetLength(length)
}

// GetData
// Get internal lump data structure
// Returns interface{} to fulfill interface
// Should be typecasted to expected type
func (lump *Visibility) GetData() *primitives.Vis {
	return &lump.data
}

// ToBytes
// Convert internal data structure into a byte slice
func (lump *Visibility) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data.NumClusters)
	for _, offset := range lump.data.ByteOffset {
		binary.Write(&buf, binary.LittleEndian, offset)
	}
	binary.Write(&buf, binary.LittleEndian, lump.data.BitVectors)
	return buf.Bytes()
}
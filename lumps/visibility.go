package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/visibility"
	"log"
)

// Lump 4: Visibility
type Visibility struct {
	LumpGeneric
	data primitives.Vis
}

// Import this lump from raw byte data
func (lump *Visibility) Unmarshall(raw []byte, length int32) {
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data.NumClusters)
	if err != nil {
		log.Fatal(err)
	}
	lump.data.ByteOffset = make([][2]int32, lump.data.NumClusters)
	err = binary.Read(bytes.NewBuffer(raw[4:]), binary.LittleEndian, &lump.data.ByteOffset)
	if err != nil {
		log.Fatal(err)
	}
	lump.data.BitVectors = make([]byte, length)
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data.BitVectors)
	if err != nil {
		log.Fatal(err)
	}

	lump.LumpInfo.SetLength(length)
}

// Get internal format structure data
func (lump *Visibility) GetData() *primitives.Vis {
	return &lump.data
}

// Dump this lump back to raw byte data
func (lump *Visibility) Marshall() ([]byte,error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data.BitVectors)
	return buf.Bytes(),err
}

package lumps

import (
	"bytes"
	"encoding/binary"
	"log"
)

// LeafFace is Lump 16: LeafFace
type LeafFace struct {
	LumpGeneric
	data []uint16 // MAX_MAP_LEAFFACES = 65536
}

// Unmarshall Imports this lump from raw byte data
func (lump *LeafFace) Unmarshall(raw []byte, length int32) {
	lump.data = make([]uint16, length/2)
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

// GetData gets internal format structure data
func (lump *LeafFace) GetData() []uint16 {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *LeafFace) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

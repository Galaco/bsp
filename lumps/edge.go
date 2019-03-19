package lumps

import (
	"bytes"
	"encoding/binary"
	"log"
)

// Edge is Lump 12: Edge
type Edge struct {
	LumpGeneric
	data [][2]uint16 // MAX_MAP_EDGES = 256000
}

// Unmarshall Imports this lump from raw byte data
func (lump *Edge) Unmarshall(raw []byte, length int32) {
	lump.data = make([][2]uint16, length/4)
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

// GetData gets internal format structure data
func (lump *Edge) GetData() [][2]uint16 {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *Edge) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

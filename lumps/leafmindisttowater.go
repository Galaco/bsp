package lumps

import (
	"bytes"
	"encoding/binary"
)

// LeafMinDistToWater is Lump 46: LeafMinDistToWater
type LeafMinDistToWater struct {
	Generic
	data []uint16
}

// Unmarshall Imports this lump from raw byte data
func (lump *LeafMinDistToWater) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([]uint16, length/2)
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return err
}

// GetData gets internal format structure data
func (lump *LeafMinDistToWater) GetData() []uint16 {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *LeafMinDistToWater) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

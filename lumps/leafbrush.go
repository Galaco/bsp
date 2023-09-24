package lumps

import (
	"bytes"
	"encoding/binary"
)

// LeafBrush is Lump 17: LeafBrush
type LeafBrush struct {
	Metadata
	data []uint16 // MAX_MAP_LEAFBRUSHES = 65536
}

// Unmarshall Imports this lump from raw byte data
func (lump *LeafBrush) Unmarshall(raw []byte) (err error) {
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
func (lump *LeafBrush) GetData() []uint16 {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *LeafBrush) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

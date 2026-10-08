package lumps

import (
	"bytes"
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
)

// Edge is Lump 12: Edge
type Edge struct {
	Generic
	data [][2]uint16 // MAX_MAP_EDGES = 256000
}

// Unmarshall Imports this lump from raw byte data
func (lump *Edge) Unmarshall(raw []byte) error {
	count, err := format.FixedRecordCount(len(raw), 4)
	if err != nil {
		return err
	}
	data := make([][2]uint16, count)
	if err := binary.Read(bytes.NewReader(raw), binary.LittleEndian, data); err != nil {
		return err
	}
	lump.data = data
	lump.Metadata.SetLength(len(raw))
	return nil
}

// GetData gets internal format structure data
func (lump *Edge) GetData() [][2]uint16 {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *Edge) Marshall() ([]byte, error) {
	length, err := format.FixedRecordBytes(len(lump.data), 4)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(make([]byte, 0, length))
	if err := binary.Write(buf, binary.LittleEndian, lump.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

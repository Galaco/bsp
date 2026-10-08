package lumps

import (
	"bytes"
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
)

// Surfedge is Lump 13: Surfedge
type Surfedge struct {
	Generic
	data []int32 // MAX_MAP_SURFEDGES = 512000
}

// Unmarshall Imports this lump from raw byte data
func (lump *Surfedge) Unmarshall(raw []byte) error {
	count, err := format.FixedRecordCount(len(raw), 4)
	if err != nil {
		return err
	}
	data := make([]int32, count)
	if err := binary.Read(bytes.NewReader(raw), binary.LittleEndian, data); err != nil {
		return err
	}
	lump.data = data
	lump.Metadata.SetLength(len(raw))
	return nil
}

// GetData gets internal format structure data
func (lump *Surfedge) GetData() []int32 {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *Surfedge) Marshall() ([]byte, error) {
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

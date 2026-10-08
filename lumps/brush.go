package lumps

import (
	"bytes"
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	primitives "github.com/galaco/bsp/primitives/brush"
)

// Brush is Lump 18: Brush
type Brush struct {
	Generic
	data []primitives.Brush
}

// Unmarshall Imports this lump from raw byte data
func (lump *Brush) Unmarshall(raw []byte) error {
	count, err := format.FixedRecordCount(len(raw), 12)
	if err != nil {
		return err
	}
	data := make([]primitives.Brush, count)
	if err := binary.Read(bytes.NewReader(raw), binary.LittleEndian, data); err != nil {
		return err
	}
	lump.data = data
	lump.Metadata.SetLength(len(raw))
	return nil
}

// GetData gets internal format structure data
func (lump *Brush) GetData() []primitives.Brush {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *Brush) Marshall() ([]byte, error) {
	length, err := format.FixedRecordBytes(len(lump.data), 12)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(make([]byte, 0, length))
	if err := binary.Write(buf, binary.LittleEndian, lump.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

package lumps

import (
	"bytes"
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	primitives "github.com/galaco/bsp/primitives/plane"
)

// Planes is Lump 1: Planes
type Planes struct {
	Generic
	data []primitives.Plane // MAP_MAX_PLANES = 65536
}

// Unmarshall Imports this lump from raw byte data
func (lump *Planes) Unmarshall(raw []byte) error {
	count, err := format.FixedRecordCount(len(raw), 20)
	if err != nil {
		return err
	}
	data := make([]primitives.Plane, count)
	if err := binary.Read(bytes.NewReader(raw), binary.LittleEndian, data); err != nil {
		return err
	}
	lump.data = data
	lump.Metadata.SetLength(len(raw))
	return nil
}

// GetData gets internal format structure data
func (lump *Planes) GetData() []primitives.Plane {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *Planes) Marshall() ([]byte, error) {
	length, err := format.FixedRecordBytes(len(lump.data), 20)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(make([]byte, 0, length))
	if err := binary.Write(buf, binary.LittleEndian, lump.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

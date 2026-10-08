package lumps

import (
	"bytes"
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	primitives "github.com/galaco/bsp/primitives/disptris"
)

// DispTris is Lump 48: DispTris
type DispTris struct {
	Generic
	data []primitives.DispTri
}

// Unmarshall Imports this lump from raw byte data
func (lump *DispTris) Unmarshall(raw []byte) error {
	count, err := format.FixedRecordCount(len(raw), 2)
	if err != nil {
		return err
	}
	data := make([]primitives.DispTri, count)
	if err := binary.Read(bytes.NewReader(raw), binary.LittleEndian, data); err != nil {
		return err
	}
	lump.data = data
	lump.Metadata.SetLength(len(raw))
	return nil
}

// GetData gets internal format structure data
func (lump *DispTris) GetData() []primitives.DispTri {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *DispTris) Marshall() ([]byte, error) {
	length, err := format.FixedRecordBytes(len(lump.data), 2)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(make([]byte, 0, length))
	if err := binary.Write(buf, binary.LittleEndian, lump.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

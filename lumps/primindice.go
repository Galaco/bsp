package lumps

import (
	"bytes"
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
)

// PrimIndice is Lump 39: PrimIndice
type PrimIndice struct {
	Generic
	data []uint16
}

// Unmarshall Imports this lump from raw byte data
func (lump *PrimIndice) Unmarshall(raw []byte) error {
	count, err := format.FixedRecordCount(len(raw), 2)
	if err != nil {
		return err
	}
	data := make([]uint16, count)
	if err := binary.Read(bytes.NewReader(raw), binary.LittleEndian, data); err != nil {
		return err
	}
	lump.data = data
	lump.Metadata.SetLength(len(raw))
	return nil
}

// GetData gets internal format structure data
func (lump *PrimIndice) GetData() []uint16 {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *PrimIndice) Marshall() ([]byte, error) {
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

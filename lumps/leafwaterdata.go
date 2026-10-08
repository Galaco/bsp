package lumps

import (
	"bytes"
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	primitives "github.com/galaco/bsp/primitives/leafwaterdata"
)

// LeafWaterData is Lump 36: leafwaterdata
type LeafWaterData struct {
	Generic
	data []primitives.LeafWaterData
}

// Unmarshall Imports this lump from raw byte data
func (lump *LeafWaterData) Unmarshall(raw []byte) error {
	count, err := format.FixedRecordCount(len(raw), 12)
	if err != nil {
		return err
	}
	data := make([]primitives.LeafWaterData, count)
	if err := binary.Read(bytes.NewReader(raw), binary.LittleEndian, data); err != nil {
		return err
	}
	lump.data = data
	lump.Metadata.SetLength(len(raw))
	return nil
}

// GetData gets internal format structure data
func (lump *LeafWaterData) GetData() []primitives.LeafWaterData {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *LeafWaterData) Marshall() ([]byte, error) {
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

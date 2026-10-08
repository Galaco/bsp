package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

// DecodeBrushes decodes complete 12-byte wire records into owned values.
func DecodeBrushes(raw []byte) ([]lumps.Brush, error) {
	count, err := format.FixedRecordCount(len(raw), 12)
	if err != nil {
		return nil, err
	}
	records := make([]lumps.Brush, count)
	for i := range records {
		wire := raw[i*12 : (i+1)*12]
		records[i].FirstSide = int32(binary.LittleEndian.Uint32(wire[0:]))
		records[i].NumSides = int32(binary.LittleEndian.Uint32(wire[4:]))
		records[i].Contents = int32(binary.LittleEndian.Uint32(wire[8:]))
	}
	return records, nil
}

// EncodeBrushes preserves every field, including named padding and float bits.
func EncodeBrushes(records []lumps.Brush) ([]byte, error) {
	size, err := format.FixedRecordBytes(len(records), 12)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*12 : (i+1)*12]
		binary.LittleEndian.PutUint32(wire[0:], uint32(records[i].FirstSide))
		binary.LittleEndian.PutUint32(wire[4:], uint32(records[i].NumSides))
		binary.LittleEndian.PutUint32(wire[8:], uint32(records[i].Contents))
	}
	return raw, nil
}

// Brushes is the private decoder adapter for fixed wire records.
type Brushes struct {
	Metadata
	data []lumps.Brush
}

func (lump *Brushes) Decode(raw []byte) error {
	data, err := DecodeBrushes(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *Brushes) Encode() ([]byte, error) { return EncodeBrushes(lump.data) }
func (lump *Brushes) GetData() []lumps.Brush  { return lump.data }

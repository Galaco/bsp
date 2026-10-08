package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

// DecodeDisplacementTriangles decodes complete 2-byte wire records into owned values.
func DecodeDisplacementTriangles(raw []byte) ([]lumps.DisplacementTriangle, error) {
	count, err := format.FixedRecordCount(len(raw), 2)
	if err != nil {
		return nil, err
	}
	records := make([]lumps.DisplacementTriangle, count)
	for i := range records {
		wire := raw[i*2 : (i+1)*2]
		records[i].Tags = binary.LittleEndian.Uint16(wire[0:])
	}
	return records, nil
}

// EncodeDisplacementTriangles preserves every field, including named padding and float bits.
func EncodeDisplacementTriangles(records []lumps.DisplacementTriangle) ([]byte, error) {
	size, err := format.FixedRecordBytes(len(records), 2)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*2 : (i+1)*2]
		binary.LittleEndian.PutUint16(wire[0:], uint16(records[i].Tags))
	}
	return raw, nil
}

// DisplacementTriangles is the private decoder adapter for fixed wire records.
type DisplacementTriangles struct {
	Metadata
	data []lumps.DisplacementTriangle
}

func (lump *DisplacementTriangles) Decode(raw []byte) error {
	data, err := DecodeDisplacementTriangles(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *DisplacementTriangles) Encode() ([]byte, error) {
	return EncodeDisplacementTriangles(lump.data)
}
func (lump *DisplacementTriangles) GetData() []lumps.DisplacementTriangle { return lump.data }

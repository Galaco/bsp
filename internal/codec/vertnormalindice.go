package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
)

// DecodeVertexNormalIndices decodes complete 2-byte wire records into owned values.
func DecodeVertexNormalIndices(raw []byte) ([]uint16, error) {
	count, err := format.FixedRecordCount(len(raw), 2)
	if err != nil {
		return nil, err
	}
	records := make([]uint16, count)
	for i := range records {
		wire := raw[i*2 : (i+1)*2]
		records[i] = binary.LittleEndian.Uint16(wire[0:])
	}
	return records, nil
}

// EncodeVertexNormalIndices preserves every field, including named padding and float bits.
func EncodeVertexNormalIndices(records []uint16) ([]byte, error) {
	size, err := format.FixedRecordBytes(len(records), 2)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*2 : (i+1)*2]
		binary.LittleEndian.PutUint16(wire[0:], uint16(records[i]))
	}
	return raw, nil
}

// VertexNormalIndices is the private decoder adapter for fixed wire records.
type VertexNormalIndices struct {
	Metadata
	data []uint16
}

func (lump *VertexNormalIndices) Decode(raw []byte) error {
	data, err := DecodeVertexNormalIndices(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *VertexNormalIndices) Encode() ([]byte, error) {
	return EncodeVertexNormalIndices(lump.data)
}
func (lump *VertexNormalIndices) GetData() []uint16 { return lump.data }

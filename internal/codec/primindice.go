package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
)

// DecodePrimitiveIndices decodes complete 2-byte wire records into owned values.
func DecodePrimitiveIndices(raw []byte) ([]uint16, error) {
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

// EncodePrimitiveIndices preserves every field, including named padding and float bits.
func EncodePrimitiveIndices(records []uint16) ([]byte, error) {
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

// PrimitiveIndices is the private decoder adapter for fixed wire records.
type PrimitiveIndices struct {
	Metadata
	data []uint16
}

func (lump *PrimitiveIndices) Decode(raw []byte) error {
	data, err := DecodePrimitiveIndices(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *PrimitiveIndices) Encode() ([]byte, error) { return EncodePrimitiveIndices(lump.data) }
func (lump *PrimitiveIndices) GetData() []uint16       { return lump.data }

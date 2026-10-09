package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
)

// DecodeEdges decodes complete 4-byte wire records into owned values.
func DecodeEdges(raw []byte) ([][2]uint16, error) {
	count, err := format.FixedRecordCount(len(raw), 4)
	if err != nil {
		return nil, err
	}
	records := make([][2]uint16, count)
	for i := range records {
		wire := raw[i*4 : (i+1)*4]
		records[i][0] = binary.LittleEndian.Uint16(wire[0:])
		records[i][1] = binary.LittleEndian.Uint16(wire[2:])
	}
	return records, nil
}

// EncodeEdges preserves every field, including named padding and float bits.
func EncodeEdges(records [][2]uint16) ([]byte, error) {
	size, err := format.FixedRecordBytes(len(records), 4)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*4 : (i+1)*4]
		binary.LittleEndian.PutUint16(wire[0:], uint16(records[i][0]))
		binary.LittleEndian.PutUint16(wire[2:], uint16(records[i][1]))
	}
	return raw, nil
}

// Edges is the private decoder adapter for fixed wire records.
type Edges struct {
	Metadata
	data [][2]uint16
}

func (lump *Edges) Decode(raw []byte) error {
	data, err := DecodeEdges(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *Edges) Encode() ([]byte, error) { return EncodeEdges(lump.data) }
func (lump *Edges) GetData() [][2]uint16    { return lump.data }

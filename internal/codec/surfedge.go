package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
)

// DecodeSurfaceEdges decodes complete 4-byte wire records into owned values.
func DecodeSurfaceEdges(raw []byte) ([]int32, error) {
	count, err := format.FixedRecordCount(len(raw), 4)
	if err != nil {
		return nil, err
	}
	records := make([]int32, count)
	for i := range records {
		wire := raw[i*4 : (i+1)*4]
		records[i] = int32(binary.LittleEndian.Uint32(wire[0:]))
	}
	return records, nil
}

// EncodeSurfaceEdges preserves every field, including named padding and float bits.
func EncodeSurfaceEdges(records []int32) ([]byte, error) {
	size, err := format.FixedRecordBytes(len(records), 4)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*4 : (i+1)*4]
		binary.LittleEndian.PutUint32(wire[0:], uint32(records[i]))
	}
	return raw, nil
}

// SurfaceEdges is the private decoder adapter for fixed wire records.
type SurfaceEdges struct {
	Metadata
	data []int32
}

func (lump *SurfaceEdges) Decode(raw []byte) error {
	data, err := DecodeSurfaceEdges(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *SurfaceEdges) Encode() ([]byte, error) { return EncodeSurfaceEdges(lump.data) }
func (lump *SurfaceEdges) GetData() []int32        { return lump.data }

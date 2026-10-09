package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
	"math"
)

// DecodeVertexNormals decodes complete 12-byte wire records into owned values.
func DecodeVertexNormals(raw []byte) ([]lumps.VertexNormal, error) {
	count, err := format.FixedRecordCount(len(raw), 12)
	if err != nil {
		return nil, err
	}
	records := make([]lumps.VertexNormal, count)
	for i := range records {
		wire := raw[i*12 : (i+1)*12]
		records[i].Pos[0] = math.Float32frombits(binary.LittleEndian.Uint32(wire[0:]))
		records[i].Pos[1] = math.Float32frombits(binary.LittleEndian.Uint32(wire[4:]))
		records[i].Pos[2] = math.Float32frombits(binary.LittleEndian.Uint32(wire[8:]))
	}
	return records, nil
}

// EncodeVertexNormals preserves every field, including named padding and float bits.
func EncodeVertexNormals(records []lumps.VertexNormal) ([]byte, error) {
	size, err := format.FixedRecordBytes(len(records), 12)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*12 : (i+1)*12]
		binary.LittleEndian.PutUint32(wire[0:], math.Float32bits(records[i].Pos[0]))
		binary.LittleEndian.PutUint32(wire[4:], math.Float32bits(records[i].Pos[1]))
		binary.LittleEndian.PutUint32(wire[8:], math.Float32bits(records[i].Pos[2]))
	}
	return raw, nil
}

// VertexNormals is the private decoder adapter for fixed wire records.
type VertexNormals struct {
	Metadata
	data []lumps.VertexNormal
}

func (lump *VertexNormals) Decode(raw []byte) error {
	data, err := DecodeVertexNormals(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *VertexNormals) Encode() ([]byte, error)       { return EncodeVertexNormals(lump.data) }
func (lump *VertexNormals) GetData() []lumps.VertexNormal { return lump.data }

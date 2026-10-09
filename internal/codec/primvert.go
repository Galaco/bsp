package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
	"math"
)

// DecodePrimitiveVertices decodes complete 12-byte wire records into owned values.
func DecodePrimitiveVertices(raw []byte) ([]lumps.PrimitiveVertex, error) {
	count, err := format.FixedRecordCount(len(raw), 12)
	if err != nil {
		return nil, err
	}
	records := make([]lumps.PrimitiveVertex, count)
	for i := range records {
		wire := raw[i*12 : (i+1)*12]
		records[i].Pos[0] = math.Float32frombits(binary.LittleEndian.Uint32(wire[0:]))
		records[i].Pos[1] = math.Float32frombits(binary.LittleEndian.Uint32(wire[4:]))
		records[i].Pos[2] = math.Float32frombits(binary.LittleEndian.Uint32(wire[8:]))
	}
	return records, nil
}

// EncodePrimitiveVertices preserves every field, including named padding and float bits.
func EncodePrimitiveVertices(records []lumps.PrimitiveVertex) ([]byte, error) {
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

// PrimitiveVertices is the private decoder adapter for fixed wire records.
type PrimitiveVertices struct {
	Metadata
	data []lumps.PrimitiveVertex
}

func (lump *PrimitiveVertices) Decode(raw []byte) error {
	data, err := DecodePrimitiveVertices(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *PrimitiveVertices) Encode() ([]byte, error)          { return EncodePrimitiveVertices(lump.data) }
func (lump *PrimitiveVertices) GetData() []lumps.PrimitiveVertex { return lump.data }

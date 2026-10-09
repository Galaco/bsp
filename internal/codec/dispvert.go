package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
	"math"
)

// DecodeDisplacementVertices decodes complete 20-byte wire records into owned values.
func DecodeDisplacementVertices(raw []byte) ([]lumps.DisplacementVertex, error) {
	count, err := format.FixedRecordCount(len(raw), 20)
	if err != nil {
		return nil, err
	}
	records := make([]lumps.DisplacementVertex, count)
	for i := range records {
		wire := raw[i*20 : (i+1)*20]
		records[i].Vec[0] = math.Float32frombits(binary.LittleEndian.Uint32(wire[0:]))
		records[i].Vec[1] = math.Float32frombits(binary.LittleEndian.Uint32(wire[4:]))
		records[i].Vec[2] = math.Float32frombits(binary.LittleEndian.Uint32(wire[8:]))
		records[i].Dist = math.Float32frombits(binary.LittleEndian.Uint32(wire[12:]))
		records[i].Alpha = math.Float32frombits(binary.LittleEndian.Uint32(wire[16:]))
	}
	return records, nil
}

// EncodeDisplacementVertices preserves every field, including named padding and float bits.
func EncodeDisplacementVertices(records []lumps.DisplacementVertex) ([]byte, error) {
	size, err := format.FixedRecordBytes(len(records), 20)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*20 : (i+1)*20]
		binary.LittleEndian.PutUint32(wire[0:], math.Float32bits(records[i].Vec[0]))
		binary.LittleEndian.PutUint32(wire[4:], math.Float32bits(records[i].Vec[1]))
		binary.LittleEndian.PutUint32(wire[8:], math.Float32bits(records[i].Vec[2]))
		binary.LittleEndian.PutUint32(wire[12:], math.Float32bits(records[i].Dist))
		binary.LittleEndian.PutUint32(wire[16:], math.Float32bits(records[i].Alpha))
	}
	return raw, nil
}

// DisplacementVertices is the private decoder adapter for fixed wire records.
type DisplacementVertices struct {
	Metadata
	data []lumps.DisplacementVertex
}

func (lump *DisplacementVertices) Decode(raw []byte) error {
	data, err := DecodeDisplacementVertices(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *DisplacementVertices) Encode() ([]byte, error) {
	return EncodeDisplacementVertices(lump.data)
}
func (lump *DisplacementVertices) GetData() []lumps.DisplacementVertex { return lump.data }

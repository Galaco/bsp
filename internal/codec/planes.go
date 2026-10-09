package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
	"math"
)

// DecodePlanes decodes complete 20-byte wire records into owned values.
func DecodePlanes(raw []byte) ([]lumps.Plane, error) {
	count, err := format.FixedRecordCount(len(raw), 20)
	if err != nil {
		return nil, err
	}
	records := make([]lumps.Plane, count)
	for i := range records {
		wire := raw[i*20 : (i+1)*20]
		records[i].Normal[0] = math.Float32frombits(binary.LittleEndian.Uint32(wire[0:]))
		records[i].Normal[1] = math.Float32frombits(binary.LittleEndian.Uint32(wire[4:]))
		records[i].Normal[2] = math.Float32frombits(binary.LittleEndian.Uint32(wire[8:]))
		records[i].Distance = math.Float32frombits(binary.LittleEndian.Uint32(wire[12:]))
		records[i].AxisType = int32(binary.LittleEndian.Uint32(wire[16:]))
	}
	return records, nil
}

// EncodePlanes preserves every field, including named padding and float bits.
func EncodePlanes(records []lumps.Plane) ([]byte, error) {
	size, err := format.FixedRecordBytes(len(records), 20)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*20 : (i+1)*20]
		binary.LittleEndian.PutUint32(wire[0:], math.Float32bits(records[i].Normal[0]))
		binary.LittleEndian.PutUint32(wire[4:], math.Float32bits(records[i].Normal[1]))
		binary.LittleEndian.PutUint32(wire[8:], math.Float32bits(records[i].Normal[2]))
		binary.LittleEndian.PutUint32(wire[12:], math.Float32bits(records[i].Distance))
		binary.LittleEndian.PutUint32(wire[16:], uint32(records[i].AxisType))
	}
	return raw, nil
}

// Planes is the private decoder adapter for fixed wire records.
type Planes struct {
	Metadata
	data []lumps.Plane
}

func (lump *Planes) Decode(raw []byte) error {
	data, err := DecodePlanes(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *Planes) Encode() ([]byte, error) { return EncodePlanes(lump.data) }
func (lump *Planes) GetData() []lumps.Plane  { return lump.data }

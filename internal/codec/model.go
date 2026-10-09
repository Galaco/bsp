package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
	"math"
)

// DecodeModels decodes complete 48-byte wire records into owned values.
func DecodeModels(raw []byte) ([]lumps.Model, error) {
	count, err := format.FixedRecordCount(len(raw), 48)
	if err != nil {
		return nil, err
	}
	records := make([]lumps.Model, count)
	for i := range records {
		wire := raw[i*48 : (i+1)*48]
		records[i].Mins[0] = math.Float32frombits(binary.LittleEndian.Uint32(wire[0:]))
		records[i].Mins[1] = math.Float32frombits(binary.LittleEndian.Uint32(wire[4:]))
		records[i].Mins[2] = math.Float32frombits(binary.LittleEndian.Uint32(wire[8:]))
		records[i].Maxs[0] = math.Float32frombits(binary.LittleEndian.Uint32(wire[12:]))
		records[i].Maxs[1] = math.Float32frombits(binary.LittleEndian.Uint32(wire[16:]))
		records[i].Maxs[2] = math.Float32frombits(binary.LittleEndian.Uint32(wire[20:]))
		records[i].Origin[0] = math.Float32frombits(binary.LittleEndian.Uint32(wire[24:]))
		records[i].Origin[1] = math.Float32frombits(binary.LittleEndian.Uint32(wire[28:]))
		records[i].Origin[2] = math.Float32frombits(binary.LittleEndian.Uint32(wire[32:]))
		records[i].HeadNode = int32(binary.LittleEndian.Uint32(wire[36:]))
		records[i].FirstFace = int32(binary.LittleEndian.Uint32(wire[40:]))
		records[i].NumFaces = int32(binary.LittleEndian.Uint32(wire[44:]))
	}
	return records, nil
}

// EncodeModels preserves every field, including named padding and float bits.
func EncodeModels(records []lumps.Model) ([]byte, error) {
	size, err := format.FixedRecordBytes(len(records), 48)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*48 : (i+1)*48]
		binary.LittleEndian.PutUint32(wire[0:], math.Float32bits(records[i].Mins[0]))
		binary.LittleEndian.PutUint32(wire[4:], math.Float32bits(records[i].Mins[1]))
		binary.LittleEndian.PutUint32(wire[8:], math.Float32bits(records[i].Mins[2]))
		binary.LittleEndian.PutUint32(wire[12:], math.Float32bits(records[i].Maxs[0]))
		binary.LittleEndian.PutUint32(wire[16:], math.Float32bits(records[i].Maxs[1]))
		binary.LittleEndian.PutUint32(wire[20:], math.Float32bits(records[i].Maxs[2]))
		binary.LittleEndian.PutUint32(wire[24:], math.Float32bits(records[i].Origin[0]))
		binary.LittleEndian.PutUint32(wire[28:], math.Float32bits(records[i].Origin[1]))
		binary.LittleEndian.PutUint32(wire[32:], math.Float32bits(records[i].Origin[2]))
		binary.LittleEndian.PutUint32(wire[36:], uint32(records[i].HeadNode))
		binary.LittleEndian.PutUint32(wire[40:], uint32(records[i].FirstFace))
		binary.LittleEndian.PutUint32(wire[44:], uint32(records[i].NumFaces))
	}
	return raw, nil
}

// Models is the private decoder adapter for fixed wire records.
type Models struct {
	Metadata
	data []lumps.Model
}

func (lump *Models) Decode(raw []byte) error {
	data, err := DecodeModels(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *Models) Encode() ([]byte, error) { return EncodeModels(lump.data) }
func (lump *Models) GetData() []lumps.Model  { return lump.data }

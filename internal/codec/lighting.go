package codec

import (
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

// DecodeLighting decodes complete 4-byte wire records into owned values.
func DecodeLighting(raw []byte) ([]lumps.ColorRGBExponent32, error) {
	count, err := format.FixedRecordCount(len(raw), 4)
	if err != nil {
		return nil, err
	}
	records := make([]lumps.ColorRGBExponent32, count)
	for i := range records {
		wire := raw[i*4 : (i+1)*4]
		records[i].R = wire[0]
		records[i].G = wire[1]
		records[i].B = wire[2]
		records[i].Exponent = int8(wire[3])
	}
	return records, nil
}

// EncodeLighting preserves every field, including named padding and float bits.
func EncodeLighting(records []lumps.ColorRGBExponent32) ([]byte, error) {
	size, err := format.FixedRecordBytes(len(records), 4)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*4 : (i+1)*4]
		wire[0] = byte(records[i].R)
		wire[1] = byte(records[i].G)
		wire[2] = byte(records[i].B)
		wire[3] = byte(records[i].Exponent)
	}
	return raw, nil
}

// Lighting is the private decoder adapter for fixed wire records.
type Lighting struct {
	Metadata
	data []lumps.ColorRGBExponent32
}

func (lump *Lighting) Decode(raw []byte) error {
	data, err := DecodeLighting(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *Lighting) Encode() ([]byte, error)             { return EncodeLighting(lump.data) }
func (lump *Lighting) GetData() []lumps.ColorRGBExponent32 { return lump.data }

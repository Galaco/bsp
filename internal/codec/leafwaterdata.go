package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
	"math"
)

// DecodeLeafWaterData decodes complete 12-byte wire records into owned values.
func DecodeLeafWaterData(raw []byte) ([]lumps.LeafWaterData, error) {
	count, err := format.FixedRecordCount(len(raw), 12)
	if err != nil {
		return nil, err
	}
	records := make([]lumps.LeafWaterData, count)
	for i := range records {
		wire := raw[i*12 : (i+1)*12]
		records[i].SurfaceZ = math.Float32frombits(binary.LittleEndian.Uint32(wire[0:]))
		records[i].MinZ = math.Float32frombits(binary.LittleEndian.Uint32(wire[4:]))
		records[i].SurfaceTexInfoID = int16(binary.LittleEndian.Uint16(wire[8:]))
		records[i].Padding[0] = wire[10]
		records[i].Padding[1] = wire[11]
	}
	return records, nil
}

// EncodeLeafWaterData preserves every field, including named padding and float bits.
func EncodeLeafWaterData(records []lumps.LeafWaterData) ([]byte, error) {
	size, err := format.FixedRecordBytes(len(records), 12)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*12 : (i+1)*12]
		binary.LittleEndian.PutUint32(wire[0:], math.Float32bits(records[i].SurfaceZ))
		binary.LittleEndian.PutUint32(wire[4:], math.Float32bits(records[i].MinZ))
		binary.LittleEndian.PutUint16(wire[8:], uint16(records[i].SurfaceTexInfoID))
		wire[10] = byte(records[i].Padding[0])
		wire[11] = byte(records[i].Padding[1])
	}
	return raw, nil
}

// LeafWaterData is the private decoder adapter for fixed wire records.
type LeafWaterData struct {
	Metadata
	data []lumps.LeafWaterData
}

func (lump *LeafWaterData) Decode(raw []byte) error {
	data, err := DecodeLeafWaterData(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *LeafWaterData) Encode() ([]byte, error)        { return EncodeLeafWaterData(lump.data) }
func (lump *LeafWaterData) GetData() []lumps.LeafWaterData { return lump.data }

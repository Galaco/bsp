package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

// DecodeVertices decodes complete 12-byte wire records into owned values.
func DecodeVertices(raw []byte) ([]mgl32.Vec3, error) {
	count, err := format.FixedRecordCount(len(raw), 12)
	if err != nil {
		return nil, err
	}
	records := make([]mgl32.Vec3, count)
	for i := range records {
		wire := raw[i*12 : (i+1)*12]
		records[i][0] = math.Float32frombits(binary.LittleEndian.Uint32(wire[0:]))
		records[i][1] = math.Float32frombits(binary.LittleEndian.Uint32(wire[4:]))
		records[i][2] = math.Float32frombits(binary.LittleEndian.Uint32(wire[8:]))
	}
	return records, nil
}

// EncodeVertices preserves every field, including named padding and float bits.
func EncodeVertices(records []mgl32.Vec3) ([]byte, error) {
	size, err := format.FixedRecordBytes(len(records), 12)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*12 : (i+1)*12]
		binary.LittleEndian.PutUint32(wire[0:], math.Float32bits(records[i][0]))
		binary.LittleEndian.PutUint32(wire[4:], math.Float32bits(records[i][1]))
		binary.LittleEndian.PutUint32(wire[8:], math.Float32bits(records[i][2]))
	}
	return raw, nil
}

// Vertices is the private decoder adapter for fixed wire records.
type Vertices struct {
	Metadata
	data []mgl32.Vec3
}

func (lump *Vertices) Decode(raw []byte) error {
	data, err := DecodeVertices(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *Vertices) Encode() ([]byte, error) { return EncodeVertices(lump.data) }
func (lump *Vertices) GetData() []mgl32.Vec3   { return lump.data }

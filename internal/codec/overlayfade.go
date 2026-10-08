package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
	"math"
)

// DecodeOverlayFades decodes complete 8-byte wire records into owned values.
func DecodeOverlayFades(raw []byte) ([]lumps.OverlayFade, error) {
	count, err := format.FixedRecordCount(len(raw), 8)
	if err != nil {
		return nil, err
	}
	records := make([]lumps.OverlayFade, count)
	for i := range records {
		wire := raw[i*8 : (i+1)*8]
		records[i].FadeDistMinSq = math.Float32frombits(binary.LittleEndian.Uint32(wire[0:]))
		records[i].FadeDistMaxSq = math.Float32frombits(binary.LittleEndian.Uint32(wire[4:]))
	}
	return records, nil
}

// EncodeOverlayFades preserves every field, including named padding and float bits.
func EncodeOverlayFades(records []lumps.OverlayFade) ([]byte, error) {
	size, err := format.FixedRecordBytes(len(records), 8)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*8 : (i+1)*8]
		binary.LittleEndian.PutUint32(wire[0:], math.Float32bits(records[i].FadeDistMinSq))
		binary.LittleEndian.PutUint32(wire[4:], math.Float32bits(records[i].FadeDistMaxSq))
	}
	return raw, nil
}

// OverlayFades is the private decoder adapter for fixed wire records.
type OverlayFades struct {
	Metadata
	data []lumps.OverlayFade
}

func (lump *OverlayFades) Decode(raw []byte) error {
	data, err := DecodeOverlayFades(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *OverlayFades) Encode() ([]byte, error)      { return EncodeOverlayFades(lump.data) }
func (lump *OverlayFades) GetData() []lumps.OverlayFade { return lump.data }

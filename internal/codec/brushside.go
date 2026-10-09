package codec

import (
	"encoding/binary"
	"fmt"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

// BrushSideHasThin selects the profile's revision-0 eight-byte dialect. L4D2
// deliberately retains the SDK signed-short bevel despite being BSP21.
func BrushSideHasThin(context format.Context) (bool, error) {
	if err := validatePCContext(context); err != nil {
		return false, err
	}
	if context.LumpVersion != 0 {
		return false, fmt.Errorf("unsupported brushside revision %d for %s", context.LumpVersion, context.Profile)
	}
	return context.Profile == "csgo" || context.Profile == "reactivedrop", nil
}

// DecodeBrushSides consumes complete eight-byte records without size-based
// dialect inference or boolean normalization of either tail field.
func DecodeBrushSides(raw []byte, context format.Context) ([]lumps.BrushSide, error) {
	hasThin, err := BrushSideHasThin(context)
	if err != nil {
		return nil, err
	}
	if len(raw)%8 != 0 {
		return nil, fmt.Errorf("brushside: %d bytes is not a multiple of 8", len(raw))
	}
	records := make([]lumps.BrushSide, len(raw)/8)
	for i := range records {
		wire := raw[i*8 : (i+1)*8]
		record := &records[i]
		record.PlaneNum = binary.LittleEndian.Uint16(wire)
		record.TexInfo = int16(binary.LittleEndian.Uint16(wire[2:]))
		record.DispInfo = int16(binary.LittleEndian.Uint16(wire[4:]))
		record.HasThin = hasThin
		if hasThin {
			record.Bevel = int16(wire[6])
			record.Thin = wire[7]
		} else {
			record.Bevel = int16(binary.LittleEndian.Uint16(wire[6:]))
		}
	}
	return records, nil
}

// EncodeBrushSides rejects dialect changes and nonrepresentable byte bevels
// before allocating output, so failed setters cannot discard Thin or sign bits.
func EncodeBrushSides(records []lumps.BrushSide, context format.Context) ([]byte, error) {
	hasThin, err := BrushSideHasThin(context)
	if err != nil {
		return nil, err
	}
	size, err := lightSidePayloadSize(len(records), 8)
	if err != nil {
		return nil, err
	}
	for i := range records {
		record := &records[i]
		if record.HasThin != hasThin {
			return nil, fmt.Errorf("brushside %d: HasThin does not match %s dialect", i, context.Profile)
		}
		if hasThin {
			if record.Bevel < 0 || record.Bevel > 255 {
				return nil, fmt.Errorf("brushside %d: bevel %d is outside byte range", i, record.Bevel)
			}
		} else if record.Thin != 0 {
			return nil, fmt.Errorf("brushside %d: thin is not representable in %s dialect", i, context.Profile)
		}
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*8 : (i+1)*8]
		record := &records[i]
		binary.LittleEndian.PutUint16(wire, record.PlaneNum)
		binary.LittleEndian.PutUint16(wire[2:], uint16(record.TexInfo))
		binary.LittleEndian.PutUint16(wire[4:], uint16(record.DispInfo))
		if hasThin {
			wire[6], wire[7] = byte(record.Bevel), record.Thin
		} else {
			binary.LittleEndian.PutUint16(wire[6:], uint16(record.Bevel))
		}
	}
	return raw, nil
}

// BrushSide wraps lump 19 for the internal registry.
type BrushSide struct {
	Metadata
	data []lumps.BrushSide
}

func (lump *BrushSide) Decode(raw []byte) error {
	data, err := DecodeBrushSides(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *BrushSide) Encode() ([]byte, error) {
	return EncodeBrushSides(lump.data, lump.Context())
}

func (lump *BrushSide) GetData() []lumps.BrushSide { return lump.data }

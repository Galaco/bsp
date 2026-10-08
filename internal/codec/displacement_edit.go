package codec

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

var ErrSharedDisplacementSampleSpan = errors.New("shared displacement sample span")

// RebaseDisplacementSampleStarts patches only following byte-offset words after
// one variable stream splice. Nil output means no displacement bytes changed.
// All offsets are checked before the single copy; other wire fields stay exact.
func RebaseDisplacementSampleStarts(raw []byte, index, start, end, oldSize int, delta int64, context format.Context) ([]byte, error) {
	if err := validateSpatialContext(context, "displacement sample offsets"); err != nil {
		return nil, err
	}
	count, err := format.FixedRecordCount(len(raw), 176)
	if err != nil {
		return nil, err
	}
	if index < 0 || index >= count || start < 0 || end < start || end > oldSize {
		return nil, fmt.Errorf("invalid displacement sample splice index %d span [%d,%d) size %d", index, start, end, oldSize)
	}
	changed := false
	for i := 0; i < count; i++ {
		if i == index {
			continue
		}
		offset := int64(int32(binary.LittleEndian.Uint32(raw[i*176+44:])))
		if offset == int64(start) {
			return nil, fmt.Errorf("%w: displacement %d start %d", ErrSharedDisplacementSampleSpan, i, start)
		}
		if offset < 0 {
			continue
		}
		if offset > int64(oldSize) {
			return nil, fmt.Errorf("displacement %d sample start %d outside %d bytes", i, offset, oldSize)
		}
		if offset < int64(end) || delta == 0 {
			continue
		}
		shifted := offset + delta
		if shifted < 0 || shifted > math.MaxInt32 {
			return nil, fmt.Errorf("displacement %d shifted sample start %d exceeds signed32", i, shifted)
		}
		changed = true
	}
	if !changed {
		return nil, nil
	}
	output := append([]byte(nil), raw...)
	for i := 0; i < count; i++ {
		if i == index {
			continue
		}
		at := i*176 + 44
		offset := int64(int32(binary.LittleEndian.Uint32(raw[at:])))
		if offset >= int64(end) {
			binary.LittleEndian.PutUint32(output[at:], uint32(offset+delta))
		}
	}
	return output, nil
}

// DisplacementMultiblendCount measures the documented marked-word ownership
// without allocating displacement snapshots solely to discard their fields.
func DisplacementMultiblendCount(raw []byte, context format.Context) (int, error) {
	if err := validateSpatialContext(context, "multiblend owners"); err != nil {
		return 0, err
	}
	count, err := format.FixedRecordCount(len(raw), 176)
	if err != nil {
		return 0, err
	}
	var total int64
	for i := 0; i < count; i++ {
		word := binary.LittleEndian.Uint32(raw[i*176+24:])
		mask := lumps.DisplacementFlagMagic | lumps.DisplacementFlagHasMultiblend
		if word&mask != mask {
			continue
		}
		power := int32(binary.LittleEndian.Uint32(raw[i*176+20:]))
		if power < 2 || power > 4 {
			return 0, fmt.Errorf("displacement %d power %d outside 2..4", i, power)
		}
		side := int64(1) << uint(power)
		total += (side + 1) * (side + 1)
		if total > math.MaxInt32/80 {
			return 0, fmt.Errorf("multiblend owner count %d exceeds BSP byte budget", total)
		}
	}
	return int(total), nil
}

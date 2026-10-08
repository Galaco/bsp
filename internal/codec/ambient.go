package codec

import (
	"encoding/binary"
	"fmt"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
	"github.com/galaco/bsp/primitives/common"
)

// Valve SDK2013 bsplib.cpp, LoadLeafAmbientLighting: leaf-v0 already supplied
// inline LDR lighting and bypasses all four separate ambient lumps. For leaf-v1,
// absence of the LDR index selects legacy cubes for BOTH channels, even when
// HDR indices exist or sample descriptors contain garbage versions.
func ambientSampleSize(context format.Context) (int, error) {
	if err := validatePCContext(context); err != nil {
		return 0, err
	}
	switch context.LeafVersion {
	case 0:
		return 0, nil
	case 1:
		if context.LegacyAmbient {
			return 24, nil
		}
		if context.LumpVersion == 0 || context.LumpVersion == 1 {
			return 28, nil
		}
		return 0, fmt.Errorf("unsupported indexed ambient sample version %d", context.LumpVersion)
	default:
		return 0, fmt.Errorf("unsupported ambient leaf version %d", context.LeafVersion)
	}
}

func ambientIndicesContext(context format.Context) error {
	if err := validatePCContext(context); err != nil {
		return err
	}
	if context.LeafVersion != 0 && context.LeafVersion != 1 {
		return fmt.Errorf("unsupported ambient leaf version %d", context.LeafVersion)
	}
	if context.LeafVersion == 1 && !context.LegacyAmbient && context.LumpVersion != 0 {
		return fmt.Errorf("unsupported ambient index version %d", context.LumpVersion)
	}
	return nil
}

func ambientOutputSize(count, stride int) (int, error) {
	// BSP descriptor lengths are signed 32-bit values. Divide the widened
	// bound first so even a native-int maximum count cannot overflow.
	limit := uint64(1<<31 - 1)
	if native := uint64(^uint(0) >> 1); native < limit {
		limit = native
	}
	if stride != 0 && uint64(count) > limit/uint64(stride) {
		return 0, fmt.Errorf("ambient output length exceeds representable BSP/native size")
	}
	return int(uint64(count) * uint64(stride)), nil
}

func decodeAmbientCube(raw []byte) common.CompressedLightCube {
	var cube common.CompressedLightCube
	for side := range cube.Color {
		offset := side * 4
		cube.Color[side] = common.ColorRGBExponent32{
			R: raw[offset], G: raw[offset+1], B: raw[offset+2], Exponent: int8(raw[offset+3]),
		}
	}
	return cube
}

func encodeAmbientCube(raw []byte, cube common.CompressedLightCube) {
	for side := range cube.Color {
		offset := side * 4
		color := cube.Color[side]
		raw[offset], raw[offset+1], raw[offset+2], raw[offset+3] = color.R, color.G, color.B, byte(color.Exponent)
	}
}

func DecodeAmbientIndices(raw []byte, context format.Context) ([]lumps.AmbientIndex, error) {
	if err := ambientIndicesContext(context); err != nil {
		return nil, err
	}
	if len(raw) != 0 && (context.LeafVersion == 0 || context.LegacyAmbient) {
		return nil, fmt.Errorf("separate ambient indices are unused by leaf-v0/legacy ambient context")
	}
	if len(raw)%4 != 0 {
		return nil, fmt.Errorf("ambient index length %d is not a multiple of 4", len(raw))
	}
	count := len(raw) / 4
	if count > lumps.MaxMapLeafs {
		return nil, fmt.Errorf("ambient index count %d exceeds %d", count, lumps.MaxMapLeafs)
	}
	data := make([]lumps.AmbientIndex, count)
	for i := range data {
		data[i].AmbientSampleCount = binary.LittleEndian.Uint16(raw[i*4 : i*4+2])
		data[i].FirstAmbientSample = binary.LittleEndian.Uint16(raw[i*4+2 : i*4+4])
	}
	return data, nil
}

func EncodeAmbientIndices(data []lumps.AmbientIndex, context format.Context) ([]byte, error) {
	if err := ambientIndicesContext(context); err != nil {
		return nil, err
	}
	if len(data) != 0 && (context.LeafVersion == 0 || context.LegacyAmbient) {
		return nil, fmt.Errorf("separate ambient indices are unused by leaf-v0/legacy ambient context")
	}
	if len(data) > lumps.MaxMapLeafs {
		return nil, fmt.Errorf("ambient index count %d exceeds %d", len(data), lumps.MaxMapLeafs)
	}
	size, err := ambientOutputSize(len(data), 4)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range data {
		binary.LittleEndian.PutUint16(raw[i*4:i*4+2], data[i].AmbientSampleCount)
		binary.LittleEndian.PutUint16(raw[i*4+2:i*4+4], data[i].FirstAmbientSample)
	}
	return raw, nil
}

func DecodeAmbientSamples(raw []byte, context format.Context) ([]lumps.AmbientSample, error) {
	stride, err := ambientSampleSize(context)
	if err != nil {
		return nil, err
	}
	if stride == 0 {
		if len(raw) != 0 {
			return nil, fmt.Errorf("leaf-v0 uses inline LDR lighting; separate ambient samples are unused")
		}
		return []lumps.AmbientSample{}, nil
	}
	if len(raw)%stride != 0 {
		return nil, fmt.Errorf("ambient sample length %d is not a multiple of %d", len(raw), stride)
	}
	data := make([]lumps.AmbientSample, len(raw)/stride)
	for i := range data {
		record := raw[i*stride : (i+1)*stride]
		data[i].Cube = decodeAmbientCube(record[:24])
		data[i].HasPosition = stride == 28
		if data[i].HasPosition {
			data[i].X, data[i].Y, data[i].Z, data[i].Pad = record[24], record[25], record[26], record[27]
		}
	}
	return data, nil
}

func EncodeAmbientSamples(data []lumps.AmbientSample, context format.Context) ([]byte, error) {
	stride, err := ambientSampleSize(context)
	if err != nil {
		return nil, err
	}
	if stride == 0 && len(data) != 0 {
		return nil, fmt.Errorf("leaf-v0 uses inline LDR lighting; separate ambient samples are unused")
	}
	positioned := stride == 28
	for i := range data {
		if data[i].HasPosition != positioned {
			return nil, fmt.Errorf("ambient sample %d position layout does not match context", i)
		}
		if !positioned && (data[i].X != 0 || data[i].Y != 0 || data[i].Z != 0 || data[i].Pad != 0) {
			return nil, fmt.Errorf("ambient sample %d has nonrepresentable position/padding in legacy layout", i)
		}
	}
	size, err := ambientOutputSize(len(data), stride)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range data {
		record := raw[i*stride : (i+1)*stride]
		encodeAmbientCube(record[:24], data[i].Cube)
		if positioned {
			record[24], record[25], record[26], record[27] = data[i].X, data[i].Y, data[i].Z, data[i].Pad
		}
	}
	return raw, nil
}

// ValidateAmbient validates one channel against its leaf/sample siblings without
// creating fallback samples or synthesizing missing channels. LumpVersion here
// is the sample descriptor version, not the index descriptor version.
func ValidateAmbient(indices []lumps.AmbientIndex, samples []lumps.AmbientSample, leafCount int, context format.Context) error {
	stride, err := ambientSampleSize(context)
	if err != nil {
		return err
	}
	if leafCount < 0 || leafCount > lumps.MaxMapLeafs {
		return fmt.Errorf("invalid ambient leaf count %d", leafCount)
	}
	if len(indices) == 0 && len(samples) == 0 {
		return nil
	}
	if stride == 0 {
		return fmt.Errorf("leaf-v0 separate ambient payloads are unused")
	}
	for i := range samples {
		if samples[i].HasPosition != (stride == 28) {
			return fmt.Errorf("ambient sample %d position layout does not match context", i)
		}
		if stride == 24 && (samples[i].X != 0 || samples[i].Y != 0 || samples[i].Z != 0 || samples[i].Pad != 0) {
			return fmt.Errorf("ambient sample %d has nonrepresentable legacy position/padding", i)
		}
	}
	if context.LegacyAmbient {
		if len(indices) != 0 {
			return fmt.Errorf("legacy ambient context ignores separate ambient indices")
		}
		if len(samples) != leafCount {
			return fmt.Errorf("legacy ambient sample count %d differs from leaf count %d", len(samples), leafCount)
		}
		return nil
	}
	if len(indices) != leafCount {
		return fmt.Errorf("ambient index count %d differs from leaf count %d", len(indices), leafCount)
	}
	for i, index := range indices {
		if index.AmbientSampleCount == 0 {
			// vrad/leaf_ambient_lighting.cpp writes a nearest-lit-leaf reference
			// for empty leaves; this is not an empty sample-array span.
			if int(index.FirstAmbientSample) >= leafCount {
				return fmt.Errorf("ambient index %d references leaf %d outside %d leaves", i, index.FirstAmbientSample, leafCount)
			}
			continue
		}
		end := uint64(index.FirstAmbientSample) + uint64(index.AmbientSampleCount)
		if end > uint64(len(samples)) {
			return fmt.Errorf("ambient index %d sample span ends at %d beyond %d samples", i, end, len(samples))
		}
	}
	return nil
}

// LeafAmbientIndex and LeafAmbientIndexHDR retain independent channel data.
type LeafAmbientIndex struct {
	lumps.Metadata
	data []lumps.AmbientIndex
}

type LeafAmbientIndexHDR struct {
	lumps.Metadata
	data []lumps.AmbientIndex
}

func (lump *LeafAmbientIndex) Unmarshall(raw []byte) error {
	data, err := DecodeAmbientIndices(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}
func (lump *LeafAmbientIndex) Marshall() ([]byte, error) {
	return EncodeAmbientIndices(lump.data, lump.Context())
}
func (lump *LeafAmbientIndex) GetData() []lumps.AmbientIndex { return lump.data }

func (lump *LeafAmbientIndexHDR) Unmarshall(raw []byte) error {
	data, err := DecodeAmbientIndices(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}
func (lump *LeafAmbientIndexHDR) Marshall() ([]byte, error) {
	return EncodeAmbientIndices(lump.data, lump.Context())
}
func (lump *LeafAmbientIndexHDR) GetData() []lumps.AmbientIndex { return lump.data }

// LeafAmbientLighting and LeafAmbientLightingHDR do not infer channel fallback.
type LeafAmbientLighting struct {
	lumps.Metadata
	data []lumps.AmbientSample
}

type LeafAmbientLightingHDR struct {
	lumps.Metadata
	data []lumps.AmbientSample
}

func (lump *LeafAmbientLighting) Unmarshall(raw []byte) error {
	data, err := DecodeAmbientSamples(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}
func (lump *LeafAmbientLighting) Marshall() ([]byte, error) {
	return EncodeAmbientSamples(lump.data, lump.Context())
}
func (lump *LeafAmbientLighting) GetData() []lumps.AmbientSample { return lump.data }

func (lump *LeafAmbientLightingHDR) Unmarshall(raw []byte) error {
	data, err := DecodeAmbientSamples(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}
func (lump *LeafAmbientLightingHDR) Marshall() ([]byte, error) {
	return EncodeAmbientSamples(lump.data, lump.Context())
}
func (lump *LeafAmbientLightingHDR) GetData() []lumps.AmbientSample { return lump.data }

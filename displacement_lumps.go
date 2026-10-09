package bsp

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/galaco/bsp/internal/codec"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

// Displacements returns complete 176-byte records with checked geometry references.
func (bsp *File) Displacements() ([]lumps.Displacement, error) {
	decoded, err := bsp.decodeLump(LumpDispInfo)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.Displacements).GetData(), nil
}

// SetDisplacements commits complete records without regenerating neighbor topology.
func (bsp *File) SetDisplacements(records []lumps.Displacement) error {
	if err := bsp.checkTypedRevision(LumpDispInfo); err != nil {
		return err
	}
	if err := bsp.validateDisplacementReferences(records); err != nil {
		return err
	}
	if err := bsp.validateDisplacementMultiblendCount(records); err != nil {
		return err
	}
	raw, err := codec.EncodeDisplacements(records, bsp.codecContext(LumpDispInfo))
	return bsp.commitOrdinary(LumpDispInfo, raw, err)
}

func (bsp *File) validateDisplacementReferences(records []lumps.Displacement) error {
	if len(records) == 0 {
		return nil
	}
	for _, id := range []LumpID{LumpFaces, LumpDispVerts, LumpDispTris} {
		if err := bsp.checkTypedRevision(id); err != nil {
			return err
		}
	}
	faces, err := bsp.logicalRecordCount(LumpFaces, 56)
	if err != nil {
		return err
	}
	vertices, err := bsp.logicalRecordCount(LumpDispVerts, 20)
	if err != nil {
		return err
	}
	triangles, err := bsp.logicalRecordCount(LumpDispTris, 2)
	if err != nil {
		return err
	}
	for i := range records {
		record := &records[i]
		if record.Power < 2 || record.Power > 4 {
			return &FormatError{"displacement references", fmt.Errorf("displacement %d invalid power %d", i, record.Power)}
		}
		side := 1 << uint(record.Power)
		vertexCount, triangleCount := (side+1)*(side+1), 2*side*side
		if int(record.MapFace) >= faces {
			return &FormatError{"displacement references", fmt.Errorf("displacement %d map face %d outside %d", i, record.MapFace, faces)}
		}
		if record.DispVertStart < 0 || int64(record.DispVertStart) > int64(vertices) || int64(vertexCount) > int64(vertices)-int64(record.DispVertStart) {
			return &FormatError{"displacement references", fmt.Errorf("displacement %d vertex span %d+%d outside %d", i, record.DispVertStart, vertexCount, vertices)}
		}
		if record.DispTriStart < 0 || int64(record.DispTriStart) > int64(triangles) || int64(triangleCount) > int64(triangles)-int64(record.DispTriStart) {
			return &FormatError{"displacement references", fmt.Errorf("displacement %d triangle span %d+%d outside %d", i, record.DispTriStart, triangleCount, triangles)}
		}
	}
	return nil
}

// DisplacementSamples reads one face-derived sample span. An absent lump 34
// returns ErrLumpNotFound; it does not fabricate samples from missing lighting.
func (bsp *File) DisplacementSamples(displacement int) ([]lumps.DisplacementSample, error) {
	_, raw, start, boundary, count, triangles, err := bsp.displacementSampleView(displacement)
	if err != nil {
		return nil, err
	}
	values, _, err := codec.DecodeDisplacementSamples(raw[start:boundary], count, triangles)
	if err != nil {
		return nil, &FormatError{fmt.Sprintf("displacement %d samples", displacement), err}
	}
	return values, nil
}

// SetDisplacementSamples preserves unclaimed bytes, shifts following byte starts
// and atomically commits sample bytes and required offsets in lump 26.
func (bsp *File) SetDisplacementSamples(displacement int, values []lumps.DisplacementSample) error {
	rawInfo, raw, start, boundary, count, triangles, err := bsp.displacementSampleView(displacement)
	if err != nil {
		return err
	}
	if len(values) != count {
		return &FormatError{"displacement samples", fmt.Errorf("got %d samples, owning face requires %d", len(values), count)}
	}
	if count == 0 {
		return nil // This face owns no bytes, even when its start is shared.
	}
	used, err := codec.MeasureDisplacementSamples(raw[start:boundary], count, triangles)
	if err != nil {
		return &FormatError{fmt.Sprintf("displacement %d samples", displacement), err}
	}
	encoded, err := codec.EncodeDisplacementSamples(values, triangles)
	if err != nil {
		return &FormatError{"encode displacement samples", err}
	}
	size := int64(len(raw)) - int64(used) + int64(len(encoded))
	if size < 0 || size > int64(1<<31-1) || uint64(size) > uint64(^uint(0)>>1) {
		return &FormatError{"displacement samples", fmt.Errorf("stored size %d exceeds BSP byte budget", size)}
	}
	info, err := codec.RebaseDisplacementSampleStarts(rawInfo, displacement, start, start+used, len(raw), int64(len(encoded))-int64(used), bsp.codecContext(LumpDispInfo))
	if errors.Is(err, codec.ErrSharedDisplacementSampleSpan) {
		return fmt.Errorf("%w: %s; replace coordinated raw payloads explicitly", ErrUnsupportedLump, err)
	}
	if err != nil {
		return &FormatError{"displacement sample offsets", err}
	}
	replacement := make([]byte, int(size))
	copy(replacement, raw[:start])
	copy(replacement[start:], encoded)
	copy(replacement[start+len(encoded):], raw[start+used:])
	// All validation/encoding succeeds before either authoritative payload changes.
	descriptor := bsp.header.Lumps[LumpDispLightmapSamplePositions]
	descriptor.FourCC = [4]byte{}
	bsp.storeOwned(LumpDispLightmapSamplePositions, descriptor, replacement)
	if info != nil {
		descriptor = bsp.header.Lumps[LumpDispInfo]
		descriptor.FourCC = [4]byte{}
		bsp.storeOwned(LumpDispInfo, descriptor, info)
	}
	return nil
}

func (bsp *File) displacementSampleView(displacement int) (rawInfo, raw []byte, start, boundary, count, triangles int, err error) {
	for _, id := range []LumpID{LumpDispInfo, LumpDispLightmapSamplePositions} {
		if err = bsp.checkTypedRevision(id); err != nil {
			return
		}
	}
	rawInfo, err = bsp.decodedLumpView(LumpDispInfo)
	if err != nil {
		return
	}
	total, countErr := format.FixedRecordCount(len(rawInfo), 176)
	if countErr != nil {
		err = &FormatError{"displacement samples", countErr}
		return
	}
	if displacement < 0 || displacement >= total {
		err = fmt.Errorf("%w: %d outside %d", ErrInvalidDisplacementIndex, displacement, total)
		return
	}
	record, decodeErr := codec.DecodeDisplacementRecord(rawInfo, displacement, bsp.codecContext(LumpDispInfo))
	if decodeErr != nil {
		err = &FormatError{"displacement sample owner", decodeErr}
		return
	}
	selected := [1]lumps.Displacement{record}
	if err = bsp.validateDisplacementReferences(selected[:]); err != nil {
		return
	}
	raw, err = bsp.decodedLumpView(LumpDispLightmapSamplePositions)
	if err != nil {
		return
	}
	if len(raw) == 0 {
		err = fmt.Errorf("%w: displacement lightmap samples", ErrLumpNotFound)
		return
	}
	if record.LightmapSampleStart < 0 || int64(record.LightmapSampleStart) > int64(len(raw)) {
		err = &FormatError{"displacement samples", fmt.Errorf("displacement %d start %d outside %d bytes", displacement, record.LightmapSampleStart, len(raw))}
		return
	}
	start, boundary = int(record.LightmapSampleStart), len(raw)
	for i := 0; i < total; i++ {
		offset := int32(binary.LittleEndian.Uint32(rawInfo[i*176+44:]))
		if offset > int32(start) && int64(offset) < int64(boundary) {
			boundary = int(offset)
		}
	}
	faces, readErr := bsp.decodedLumpView(LumpFaces)
	if readErr != nil {
		err = readErr
		return
	}
	face := faces[int(record.MapFace)*56 : (int(record.MapFace)+1)*56]
	width, height := int64(int32(binary.LittleEndian.Uint32(face[36:])))+1, int64(int32(binary.LittleEndian.Uint32(face[40:])))+1
	if width < 0 || height < 0 || height != 0 && width > int64(boundary-start)/4/height {
		err = &FormatError{"displacement samples", fmt.Errorf("displacement %d dimensions %d×%d exceed byte span %d", displacement, width, height, boundary-start)}
		return
	}
	count = int(width * height)
	side := 1 << uint(record.Power)
	triangles = 2 * side * side
	return
}

// DisplacementMultiblends returns explicit weights, alpha and four RGB colors.
func (bsp *File) DisplacementMultiblends() ([]lumps.DisplacementMultiblend, error) {
	decoded, err := bsp.decodeLump(LumpDispMultiBlend)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.DisplacementMultiblends).GetData(), nil
}

// SetDisplacementMultiblends commits values only for a verified profile.
func (bsp *File) SetDisplacementMultiblends(records []lumps.DisplacementMultiblend) error {
	if err := bsp.checkTypedRevision(LumpDispMultiBlend); err != nil {
		return err
	}
	if err := bsp.validateMultiblendReferences(records); err != nil {
		return err
	}
	raw, err := codec.EncodeDisplacementMultiblends(records, bsp.codecContext(LumpDispMultiBlend))
	return bsp.commitOrdinary(LumpDispMultiBlend, raw, err)
}
func (bsp *File) validateMultiblendReferences(values []lumps.DisplacementMultiblend) error {
	if err := bsp.checkTypedRevision(LumpDispInfo); err != nil {
		return err
	}
	raw, err := bsp.decodedLumpView(LumpDispInfo)
	if err != nil {
		return err
	}
	total, err := codec.DisplacementMultiblendCount(raw, bsp.codecContext(LumpDispInfo))
	if err != nil {
		return &FormatError{"displacement multiblend owners", err}
	}
	if len(values) != total {
		return &FormatError{"displacement multiblend", fmt.Errorf("got %d records, flagged displacements require %d", len(values), total)}
	}
	return nil
}

func (bsp *File) validateDisplacementMultiblendCount(records []lumps.Displacement) error {
	if bsp.profile != ProfileCSGO && bsp.profile != ProfileReactiveDrop {
		return nil
	}
	var expected int64
	for i := range records {
		record := &records[i]
		mask := lumps.DisplacementFlagMagic | lumps.DisplacementFlagHasMultiblend
		if record.Flags&mask != mask {
			continue
		}
		side := int64(1) << uint(record.Power)
		expected += (side + 1) * (side + 1)
		if expected > int64(1<<31-1)/80 {
			return &FormatError{"displacement multiblend owners", fmt.Errorf("count %d exceeds byte budget", expected)}
		}
	}
	// An unrelated opaque revision must survive ordinary displacement edits.
	// Explicit ownership cannot be checked against a decoder we do not possess.
	if bsp.header.Lumps[LumpDispMultiBlend].Version != 0 {
		if expected == 0 {
			return nil
		}
		return fmt.Errorf("%w: multiblend owner references unknown revision %d", ErrUnsupportedLump, bsp.header.Lumps[LumpDispMultiBlend].Version)
	}
	count, err := bsp.logicalRecordCount(LumpDispMultiBlend, 80)
	if err != nil {
		return err
	}
	if int64(count) != expected {
		return &FormatError{"displacement multiblend owners", fmt.Errorf("new owners require %d records, stored channel has %d", expected, count)}
	}
	return nil
}

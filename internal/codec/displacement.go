package codec

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

// Valve b8cfb12 and Reactive Drop af3be02, utils/vbsp/disp_vbsp.cpp,
// EmitInitialDispInfos: minTess = flags | 0x80000000. The marker, not the
// profile or a guessed contents mask, selects the interpretation of byte 24.
func displacementTessWord(record *lumps.Displacement) (uint32, error) {
	if record.Flags != 0 {
		if record.Flags&lumps.DisplacementFlagMagic == 0 {
			return 0, fmt.Errorf("flags %#x lack the MinTess-word marker", record.Flags)
		}
		if record.MinTess != 0 {
			return 0, fmt.Errorf("MinTess %d conflicts with marked flags %#x", record.MinTess, record.Flags)
		}
		return record.Flags, nil
	}
	if record.MinTess < 0 {
		return 0, fmt.Errorf("MinTess %#x uses the flags marker; supply Flags instead", uint32(record.MinTess))
	}
	return uint32(record.MinTess), nil
}

func validateDisplacement(record *lumps.Displacement, index, count int) error {
	if record.Power < 2 || record.Power > 4 {
		return fmt.Errorf("displacement %d: power %d outside 2..4", index, record.Power)
	}
	if _, err := displacementTessWord(record); err != nil {
		return fmt.Errorf("displacement %d: %w", index, err)
	}
	// Check before shifting, then widen the signed starts before adding spans.
	side := int64(1) << uint(record.Power)
	vertices, triangles := (side+1)*(side+1), 2*side*side
	if record.DispVertStart < 0 || int64(record.DispVertStart)+vertices > math.MaxInt32 {
		return fmt.Errorf("displacement %d: vertex span start %d count %d exceeds signed-32-bit range", index, record.DispVertStart, vertices)
	}
	if record.DispTriStart < 0 || int64(record.DispTriStart)+triangles > math.MaxInt32 {
		return fmt.Errorf("displacement %d: triangle span start %d count %d exceeds signed-32-bit range", index, record.DispTriStart, triangles)
	}
	for edge := range record.EdgeNeighbors {
		for sub := range record.EdgeNeighbors[edge].SubNeighbors {
			neighbor := &record.EdgeNeighbors[edge].SubNeighbors[sub]
			if neighbor.Index == 0xffff {
				continue
			}
			if int(neighbor.Index) >= count {
				return fmt.Errorf("displacement %d: edge %d subneighbor %d index %d outside %d displacements", index, edge, sub, neighbor.Index, count)
			}
			if neighbor.Orientation > 3 || neighbor.Span > 2 || neighbor.NeighborSpan > 2 {
				return fmt.Errorf("displacement %d: edge %d subneighbor %d invalid orientation/span/neighbor span %d/%d/%d", index, edge, sub, neighbor.Orientation, neighbor.Span, neighbor.NeighborSpan)
			}
		}
	}
	for corner := range record.CornerNeighbors {
		neighbor := &record.CornerNeighbors[corner]
		if neighbor.Count > 4 {
			return fmt.Errorf("displacement %d: corner %d count %d exceeds 4", index, corner, neighbor.Count)
		}
		for slot := 0; slot < int(neighbor.Count); slot++ {
			if n := neighbor.Neighbors[slot]; n != 0xffff && int(n) >= count {
				return fmt.Errorf("displacement %d: corner %d neighbor %d index %d outside %d displacements", index, corner, slot, n, count)
			}
		}
	}
	return nil
}

// DecodeDisplacements consumes complete 176-byte records. References to sibling
// face/vertex/triangle lumps are checked by the attached File, not this codec.
func DecodeDisplacements(raw []byte, context format.Context) ([]lumps.Displacement, error) {
	if err := validateSpatialContext(context, "displacement"); err != nil {
		return nil, err
	}
	count, err := format.FixedRecordCount(len(raw), 176)
	if err != nil {
		return nil, fmt.Errorf("displacement: %w", err)
	}
	records := make([]lumps.Displacement, count)
	for i := range records {
		decodeDisplacementWire(raw[i*176:(i+1)*176], &records[i])
		if err := validateDisplacement(&records[i], i, count); err != nil {
			return nil, err
		}
	}
	return records, nil
}

// DecodeDisplacementRecord reads only the selected record, while checking its
// neighbor indices against the full displacement array. Unselected contents
// are not interpreted, but the whole payload must consist of 176-byte records.
func DecodeDisplacementRecord(raw []byte, index int, context format.Context) (lumps.Displacement, error) {
	var record lumps.Displacement
	if err := validateSpatialContext(context, "displacement"); err != nil {
		return record, err
	}
	count, err := format.FixedRecordCount(len(raw), 176)
	if err != nil {
		return record, fmt.Errorf("displacement: %w", err)
	}
	if index < 0 || index >= count {
		return record, fmt.Errorf("displacement index %d outside record count %d", index, count)
	}
	decodeDisplacementWire(raw[index*176:(index+1)*176], &record)
	if err := validateDisplacement(&record, index, count); err != nil {
		return lumps.Displacement{}, err
	}
	return record, nil
}

func decodeDisplacementWire(wire []byte, record *lumps.Displacement) {
	for j := range record.StartPosition {
		record.StartPosition[j] = math.Float32frombits(binary.LittleEndian.Uint32(wire[j*4:]))
	}
	record.DispVertStart = int32(binary.LittleEndian.Uint32(wire[12:]))
	record.DispTriStart = int32(binary.LittleEndian.Uint32(wire[16:]))
	record.Power = int32(binary.LittleEndian.Uint32(wire[20:]))
	if word := binary.LittleEndian.Uint32(wire[24:]); word&lumps.DisplacementFlagMagic != 0 {
		record.Flags = word
	} else {
		record.MinTess = int32(word)
	}
	record.SmoothingAngle = math.Float32frombits(binary.LittleEndian.Uint32(wire[28:]))
	record.Contents = int32(binary.LittleEndian.Uint32(wire[32:]))
	record.MapFace = binary.LittleEndian.Uint16(wire[36:])
	copy(record.MapFacePadding[:], wire[38:40])
	record.LightmapAlphaStart = int32(binary.LittleEndian.Uint32(wire[40:]))
	record.LightmapSampleStart = int32(binary.LittleEndian.Uint32(wire[44:]))
	for edge := range record.EdgeNeighbors {
		for sub := range record.EdgeNeighbors[edge].SubNeighbors {
			offset := 48 + edge*12 + sub*6
			neighbor := &record.EdgeNeighbors[edge].SubNeighbors[sub]
			neighbor.Index = binary.LittleEndian.Uint16(wire[offset:])
			neighbor.Orientation, neighbor.Span, neighbor.NeighborSpan, neighbor.Padding = wire[offset+2], wire[offset+3], wire[offset+4], wire[offset+5]
		}
	}
	for corner := range record.CornerNeighbors {
		offset := 96 + corner*10
		neighbor := &record.CornerNeighbors[corner]
		for slot := range neighbor.Neighbors {
			neighbor.Neighbors[slot] = binary.LittleEndian.Uint16(wire[offset+slot*2:])
		}
		neighbor.Count, neighbor.Padding = wire[offset+8], wire[offset+9]
	}
	for word := range record.AllowedVerts {
		record.AllowedVerts[word] = binary.LittleEndian.Uint32(wire[136+word*4:])
	}
}

func EncodeDisplacements(records []lumps.Displacement, context format.Context) ([]byte, error) {
	if err := validateSpatialContext(context, "displacement"); err != nil {
		return nil, err
	}
	size, err := format.FixedRecordBytes(len(records), 176)
	if err != nil {
		return nil, fmt.Errorf("displacement: %w", err)
	}
	for i := range records {
		if err := validateDisplacement(&records[i], i, len(records)); err != nil {
			return nil, err
		}
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*176 : (i+1)*176]
		record := &records[i]
		for j := range record.StartPosition {
			binary.LittleEndian.PutUint32(wire[j*4:], math.Float32bits(record.StartPosition[j]))
		}
		binary.LittleEndian.PutUint32(wire[12:], uint32(record.DispVertStart))
		binary.LittleEndian.PutUint32(wire[16:], uint32(record.DispTriStart))
		binary.LittleEndian.PutUint32(wire[20:], uint32(record.Power))
		word := uint32(record.MinTess)
		if record.Flags != 0 {
			word = record.Flags
		}
		binary.LittleEndian.PutUint32(wire[24:], word)
		binary.LittleEndian.PutUint32(wire[28:], math.Float32bits(record.SmoothingAngle))
		binary.LittleEndian.PutUint32(wire[32:], uint32(record.Contents))
		binary.LittleEndian.PutUint16(wire[36:], record.MapFace)
		copy(wire[38:40], record.MapFacePadding[:])
		binary.LittleEndian.PutUint32(wire[40:], uint32(record.LightmapAlphaStart))
		binary.LittleEndian.PutUint32(wire[44:], uint32(record.LightmapSampleStart))
		for edge := range record.EdgeNeighbors {
			for sub := range record.EdgeNeighbors[edge].SubNeighbors {
				offset := 48 + edge*12 + sub*6
				neighbor := &record.EdgeNeighbors[edge].SubNeighbors[sub]
				binary.LittleEndian.PutUint16(wire[offset:], neighbor.Index)
				wire[offset+2], wire[offset+3], wire[offset+4], wire[offset+5] = neighbor.Orientation, neighbor.Span, neighbor.NeighborSpan, neighbor.Padding
			}
		}
		for corner := range record.CornerNeighbors {
			offset := 96 + corner*10
			neighbor := &record.CornerNeighbors[corner]
			for slot := range neighbor.Neighbors {
				binary.LittleEndian.PutUint16(wire[offset+slot*2:], neighbor.Neighbors[slot])
			}
			wire[offset+8], wire[offset+9] = neighbor.Count, neighbor.Padding
		}
		for word := range record.AllowedVerts {
			binary.LittleEndian.PutUint32(wire[136+word*4:], record.AllowedVerts[word])
		}
	}
	return raw, nil
}

// DecodeDisplacementSamples reads exactly count samples and returns their byte
// consumption; following bytes belong to the owner and are never discarded.
// Primary CalculateLightmapSamplePositions uses (face LM width+1)*(height+1),
// independently of light styles. The owner supplies that count and local tris.
func DecodeDisplacementSamples(raw []byte, count, triangleCount int) ([]lumps.DisplacementSample, int, error) {
	if err := validateDisplacementSampleInput(raw, count, triangleCount); err != nil {
		return nil, 0, err
	}
	records := make([]lumps.DisplacementSample, count)
	used, err := scanDisplacementSamples(raw, count, triangleCount, records)
	if err != nil {
		return nil, used, err
	}
	return records, used, nil
}

// MeasureDisplacementSamples validates exactly count records without allocating
// sample values. The returned size excludes following, unclaimed stream bytes.
func MeasureDisplacementSamples(raw []byte, count, triangleCount int) (int, error) {
	if err := validateDisplacementSampleInput(raw, count, triangleCount); err != nil {
		return 0, err
	}
	return scanDisplacementSamples(raw, count, triangleCount, nil)
}

func validateDisplacementSampleInput(raw []byte, count, triangleCount int) error {
	if count < 0 {
		return fmt.Errorf("displacement samples: negative count %d", count)
	}
	if triangleCount < 0 {
		return fmt.Errorf("displacement samples: negative triangle count %d", triangleCount)
	}
	if count > len(raw)/4 {
		return fmt.Errorf("displacement samples: count %d exceeds %d-byte minimum sample capacity %d", count, len(raw), len(raw)/4)
	}
	if _, err := format.FixedRecordBytes(count, 4); err != nil {
		return fmt.Errorf("displacement samples: %w", err)
	}
	return nil
}

func scanDisplacementSamples(raw []byte, count, triangleCount int, records []lumps.DisplacementSample) (int, error) {
	used := 0
	for i := 0; i < count; i++ {
		if used == len(raw) {
			return used, fmt.Errorf("displacement sample %d at byte %d: missing triangle index", i, used)
		}
		index := uint16(raw[used])
		used++
		if index == 255 {
			if used == len(raw) {
				return used, fmt.Errorf("displacement sample %d at byte %d: truncated index escape", i, used)
			}
			index += uint16(raw[used])
			used++
		}
		if len(raw)-used < 3 {
			return used, fmt.Errorf("displacement sample %d at byte %d: truncated barycentric coordinates", i, used)
		}
		if int(index) >= triangleCount {
			return used, fmt.Errorf("displacement sample %d: triangle %d outside local triangle count %d", i, index, triangleCount)
		}
		if records != nil {
			records[i].Index = index
			copy(records[i].Barycentric[:], raw[used:used+3])
		}
		used += 3
	}
	return used, nil
}

func displacementSampleBytes(count, escaped int) (int, error) {
	size, err := format.FixedRecordBytes(count, 4)
	if err != nil {
		return 0, fmt.Errorf("displacement samples: %w", err)
	}
	if escaped < 0 || escaped > count {
		return 0, fmt.Errorf("displacement samples: invalid escaped count %d for %d samples", escaped, count)
	}
	limit := uint64(math.MaxInt32)
	if addressLimit := uint64(^uint(0) >> 1); addressLimit < limit {
		limit = addressLimit
	}
	if uint64(size)+uint64(escaped) > limit {
		return 0, fmt.Errorf("displacement samples: output size exceeds byte budget %d", limit)
	}
	return size + escaped, nil
}

func EncodeDisplacementSamples(records []lumps.DisplacementSample, triangleCount int) ([]byte, error) {
	if triangleCount < 0 {
		return nil, fmt.Errorf("displacement samples: negative triangle count %d", triangleCount)
	}
	if _, err := displacementSampleBytes(len(records), 0); err != nil {
		return nil, err
	}
	escaped := 0
	for i := range records {
		index := records[i].Index
		if index > 510 {
			return nil, fmt.Errorf("displacement sample %d: triangle index %d exceeds stream maximum 510", i, index)
		}
		if int(index) >= triangleCount {
			return nil, fmt.Errorf("displacement sample %d: triangle %d outside local triangle count %d", i, index, triangleCount)
		}
		if index >= 255 {
			escaped++
		}
	}
	size, err := displacementSampleBytes(len(records), escaped)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	used := 0
	for i := range records {
		if index := records[i].Index; index < 255 {
			raw[used] = byte(index)
			used++
		} else {
			raw[used], raw[used+1] = 255, byte(index-255)
			used += 2
		}
		copy(raw[used:used+3], records[i].Barycentric[:])
		used += 3
	}
	return raw, nil
}

// Displacements is the private registry adapter for the complete value codec.
type Displacements struct {
	lumps.Metadata
	data []lumps.Displacement
}

func (lump *Displacements) Unmarshall(raw []byte) error {
	data, err := DecodeDisplacements(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *Displacements) Marshall() ([]byte, error) {
	return EncodeDisplacements(lump.data, lump.Context())
}

func (lump *Displacements) GetData() []lumps.Displacement { return lump.data }

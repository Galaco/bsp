package codec

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

// Occlusion is the version-aware wrapper for lump 9.
type Occlusion struct {
	lumps.Metadata
	data lumps.Occlusion
}

func (lump *Occlusion) Unmarshall(raw []byte) error {
	value, err := DecodeOcclusion(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = value
	lump.SetLength(len(raw))
	return nil
}

func (lump *Occlusion) Marshall() ([]byte, error) {
	return EncodeOcclusion(lump.data, lump.Context())
}

func (lump *Occlusion) GetData() lumps.Occlusion { return lump.data }

// DecodeOcclusion reads three signed counted arrays, preserving all records.
// Widths follow doccluderdataV1_t/doccluderdata_t and doccluderpolydata_t in
// Valve SDK2013 bspfile.h, revision b8cfb12c0e083a2ef5b2f9f9b50f3902fa034474.
func DecodeOcclusion(raw []byte, context format.Context) (lumps.Occlusion, error) {
	width, err := occlusionRecordWidth(context)
	if err != nil {
		return lumps.Occlusion{}, err
	}
	// Locate all arrays and prove every count against its available bytes before
	// allocating. Later count words are reserved while checking earlier arrays.
	offset := 0
	dataCount, err := readOcclusionCount(raw, &offset, width, 8, "occluders")
	if err != nil {
		return lumps.Occlusion{}, err
	}
	dataOffset := offset
	offset += dataCount * width
	polyCount, err := readOcclusionCount(raw, &offset, 12, 4, "polygons")
	if err != nil {
		return lumps.Occlusion{}, err
	}
	polyOffset := offset
	offset += polyCount * 12
	indexCount, err := readOcclusionCount(raw, &offset, 4, 0, "vertex indices")
	if err != nil {
		return lumps.Occlusion{}, err
	}
	indexOffset := offset
	offset += indexCount * 4
	if offset != len(raw) {
		return lumps.Occlusion{}, fmt.Errorf("occlusion: %d trailing bytes", len(raw)-offset)
	}
	value := lumps.Occlusion{
		Data:          make([]lumps.OcclusionData, dataCount),
		Polygons:      make([]lumps.OcclusionPolyData, polyCount),
		VertexIndices: make([]int32, indexCount),
	}
	for i := range value.Data {
		record := raw[dataOffset+i*width : dataOffset+(i+1)*width]
		data := &value.Data[i]
		data.Flags = int32(binary.LittleEndian.Uint32(record[0:4]))
		data.FirstPoly = int32(binary.LittleEndian.Uint32(record[4:8]))
		data.PolyCount = int32(binary.LittleEndian.Uint32(record[8:12]))
		for axis := 0; axis < 3; axis++ {
			data.Mins[axis] = math.Float32frombits(binary.LittleEndian.Uint32(record[12+axis*4 : 16+axis*4]))
			data.Maxs[axis] = math.Float32frombits(binary.LittleEndian.Uint32(record[24+axis*4 : 28+axis*4]))
		}
		if width == 40 {
			data.Area = int32(binary.LittleEndian.Uint32(record[36:40]))
			data.HasArea = true
		}
	}
	for i := range value.Polygons {
		record := raw[polyOffset+i*12 : polyOffset+(i+1)*12]
		value.Polygons[i] = lumps.OcclusionPolyData{
			FirstVertexIndex: int32(binary.LittleEndian.Uint32(record[0:4])),
			VertexCount:      int32(binary.LittleEndian.Uint32(record[4:8])),
			PlaneNum:         int32(binary.LittleEndian.Uint32(record[8:12])),
		}
	}
	for i := range value.VertexIndices {
		value.VertexIndices[i] = int32(binary.LittleEndian.Uint32(raw[indexOffset+i*4 : indexOffset+(i+1)*4]))
	}
	if err := validateOcclusionSpans(value); err != nil {
		return lumps.Occlusion{}, err
	}
	return value, nil
}

// EncodeOcclusion derives signed wire counts from slices and rejects fields
// that cannot be represented by the selected revision before allocating output.
func EncodeOcclusion(value lumps.Occlusion, context format.Context) ([]byte, error) {
	width, err := occlusionRecordWidth(context)
	if err != nil {
		return nil, err
	}
	size, err := occlusionEncodedSize(len(value.Data), len(value.Polygons), len(value.VertexIndices), width)
	if err != nil {
		return nil, err
	}
	if err := validateOcclusionSpans(value); err != nil {
		return nil, err
	}
	if width == 36 {
		for i, data := range value.Data {
			if data.HasArea || data.Area != 0 {
				return nil, fmt.Errorf("occlusion: occluder %d area is not representable in revision 1", i)
			}
		}
	}
	raw := make([]byte, size)
	binary.LittleEndian.PutUint32(raw[0:4], uint32(len(value.Data)))
	offset := 4
	for _, data := range value.Data {
		record := raw[offset : offset+width]
		binary.LittleEndian.PutUint32(record[0:4], uint32(data.Flags))
		binary.LittleEndian.PutUint32(record[4:8], uint32(data.FirstPoly))
		binary.LittleEndian.PutUint32(record[8:12], uint32(data.PolyCount))
		for axis := 0; axis < 3; axis++ {
			binary.LittleEndian.PutUint32(record[12+axis*4:16+axis*4], math.Float32bits(data.Mins[axis]))
			binary.LittleEndian.PutUint32(record[24+axis*4:28+axis*4], math.Float32bits(data.Maxs[axis]))
		}
		if width == 40 {
			binary.LittleEndian.PutUint32(record[36:40], uint32(data.Area))
		}
		offset += width
	}
	binary.LittleEndian.PutUint32(raw[offset:offset+4], uint32(len(value.Polygons)))
	offset += 4
	for _, poly := range value.Polygons {
		binary.LittleEndian.PutUint32(raw[offset:offset+4], uint32(poly.FirstVertexIndex))
		binary.LittleEndian.PutUint32(raw[offset+4:offset+8], uint32(poly.VertexCount))
		binary.LittleEndian.PutUint32(raw[offset+8:offset+12], uint32(poly.PlaneNum))
		offset += 12
	}
	binary.LittleEndian.PutUint32(raw[offset:offset+4], uint32(len(value.VertexIndices)))
	offset += 4
	for _, index := range value.VertexIndices {
		binary.LittleEndian.PutUint32(raw[offset:offset+4], uint32(index))
		offset += 4
	}
	return raw, nil
}

func occlusionRecordWidth(context format.Context) (int, error) {
	if err := validatePCContext(context); err != nil {
		return 0, fmt.Errorf("occlusion: %w", err)
	}
	switch context.LumpVersion {
	case 1:
		return 36, nil
	case 2:
		return 40, nil
	default:
		return 0, fmt.Errorf("occlusion: unsupported revision %d", context.LumpVersion)
	}
}

func readOcclusionCount(raw []byte, offset *int, width, reserved int, name string) (int, error) {
	if len(raw)-*offset < 4 {
		return 0, fmt.Errorf("occlusion: %s count: %w", name, io.ErrUnexpectedEOF)
	}
	count := int32(binary.LittleEndian.Uint32(raw[*offset : *offset+4]))
	*offset += 4
	if count < 0 {
		return 0, fmt.Errorf("occlusion: negative %s count %d", name, count)
	}
	remaining := len(raw) - *offset
	if remaining < reserved || int64(count) > int64((remaining-reserved)/width) {
		return 0, fmt.Errorf("occlusion: %s count %d exceeds remaining payload: %w", name, count, io.ErrUnexpectedEOF)
	}
	return int(count), nil
}

func occlusionEncodedSize(dataCount, polyCount, indexCount, width int) (int, error) {
	const maxCount = int64(1<<31 - 1)
	counts := [3]int{dataCount, polyCount, indexCount}
	widths := [3]int{width, 12, 4}
	size := int64(12)
	for i, count := range counts {
		if count < 0 || int64(count) > maxCount {
			return 0, fmt.Errorf("occlusion: array %d count %d is outside signed 32-bit range", i, count)
		}
		bytes := int64(count) * int64(widths[i])
		if bytes > maxCount-size {
			return 0, fmt.Errorf("occlusion: encoded payload exceeds signed 32-bit size")
		}
		size += bytes
	}
	if uint64(size) > uint64(^uint(0)>>1) {
		return 0, fmt.Errorf("occlusion: encoded payload exceeds native int size")
	}
	return int(size), nil
}

func validateOcclusionSpans(value lumps.Occlusion) error {
	for i, data := range value.Data {
		if !validOcclusionSpan(data.FirstPoly, data.PolyCount, len(value.Polygons)) {
			return fmt.Errorf("occlusion: occluder %d polygon span (%d, %d) exceeds %d polygons", i, data.FirstPoly, data.PolyCount, len(value.Polygons))
		}
	}
	for i, poly := range value.Polygons {
		if !validOcclusionSpan(poly.FirstVertexIndex, poly.VertexCount, len(value.VertexIndices)) {
			return fmt.Errorf("occlusion: polygon %d vertex span (%d, %d) exceeds %d indices", i, poly.FirstVertexIndex, poly.VertexCount, len(value.VertexIndices))
		}
	}
	return nil
}

func validOcclusionSpan(first, count int32, length int) bool {
	return first >= 0 && count >= 0 && int64(first)+int64(count) <= int64(length)
}

// ValidateOcclusionReferences validates self-contained spans and every external
// plane and vertex reference, including records not selected by an occluder.
// Bounds are exclusive; no negative reference is a sentinel in this payload.
func ValidateOcclusionReferences(value lumps.Occlusion, planeCount, vertexCount int) error {
	if planeCount < 0 || vertexCount < 0 {
		return fmt.Errorf("occlusion: negative plane or vertex count")
	}
	if err := validateOcclusionSpans(value); err != nil {
		return err
	}
	for i, poly := range value.Polygons {
		if poly.PlaneNum < 0 || int64(poly.PlaneNum) >= int64(planeCount) {
			return fmt.Errorf("occlusion: polygon %d plane %d outside %d planes", i, poly.PlaneNum, planeCount)
		}
	}
	for i, index := range value.VertexIndices {
		if index < 0 || int64(index) >= int64(vertexCount) {
			return fmt.Errorf("occlusion: vertex index %d references %d outside %d vertices", i, index, vertexCount)
		}
	}
	return nil
}

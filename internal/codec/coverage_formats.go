package codec

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
	"github.com/go-gl/mathgl/mgl32"
)

// Widths and field order follow the pinned Reactive Drop bspfile.h and its
// bsplib.cpp CopyLump/AddLump consumers (af3be02c1d7db2cd07c41105e7519cf8edd66ce3).
// Candidate RD wire contexts do not themselves admit a root File profile.
func validateCoverageContext(context format.Context, name string) error {
	if context.Profile == "reactivedrop" {
		if context.ByteOrder != nil && context.ByteOrder != binary.LittleEndian {
			return fmt.Errorf("%s: unsupported byte order", name)
		}
		if context.BSPVersion < 19 || context.BSPVersion > 21 {
			return fmt.Errorf("%s: unsupported Reactive Drop BSP version %d", name, context.BSPVersion)
		}
	} else if err := validatePCContext(context); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	if context.LumpVersion != 0 {
		return fmt.Errorf("%s: unsupported revision %d", name, context.LumpVersion)
	}
	return nil
}

func validateCoverageModernContext(context format.Context, name string, allowL4D2 bool) error {
	if err := validateCoverageContext(context, name); err != nil {
		return err
	}
	if context.Profile == "csgo" || context.Profile == "reactivedrop" || (allowL4D2 && context.Profile == "l4d2") {
		return nil
	}
	return fmt.Errorf("%s: unsupported profile %q", name, context.Profile)
}

func validatePropCoverageContext(context format.Context, name string) error {
	if context.Profile != "reactivedrop" {
		return fmt.Errorf("%s: unsupported profile %q", name, context.Profile)
	}
	return validateCoverageContext(context, name)
}

func coverageVec3(raw []byte) mgl32.Vec3 {
	return mgl32.Vec3{
		math.Float32frombits(binary.LittleEndian.Uint32(raw)),
		math.Float32frombits(binary.LittleEndian.Uint32(raw[4:])),
		math.Float32frombits(binary.LittleEndian.Uint32(raw[8:])),
	}
}

func putCoverageVec3(raw []byte, value mgl32.Vec3) {
	for j := range value {
		binary.LittleEndian.PutUint32(raw[j*4:], math.Float32bits(value[j]))
	}
}

// DecodeWaterOverlays retains all 256 face slots and complete geometry. The
// attached File, rather than this wire decoder, validates sibling references.
func DecodeWaterOverlays(raw []byte, context format.Context) ([]lumps.WaterOverlay, error) {
	if err := validateCoverageContext(context, "water overlay"); err != nil {
		return nil, err
	}
	count, err := format.FixedRecordCount(len(raw), 1120)
	if err != nil {
		return nil, fmt.Errorf("water overlay: %w", err)
	}
	for i := 0; i < count; i++ {
		if faces := binary.LittleEndian.Uint16(raw[i*1120+6:]) & 0x3fff; faces > 256 {
			return nil, fmt.Errorf("water overlay %d: face count %d exceeds 256", i, faces)
		}
	}
	records := make([]lumps.WaterOverlay, count)
	for i := range records {
		wire := raw[i*1120 : (i+1)*1120]
		record := &records[i]
		record.ID = int32(binary.LittleEndian.Uint32(wire))
		record.TexInfo = int16(binary.LittleEndian.Uint16(wire[4:]))
		record.FaceCountAndRenderOrder = binary.LittleEndian.Uint16(wire[6:])
		for j := range record.Faces {
			record.Faces[j] = int32(binary.LittleEndian.Uint32(wire[8+j*4:]))
		}
		for j := 0; j < 2; j++ {
			record.U[j] = math.Float32frombits(binary.LittleEndian.Uint32(wire[1032+j*4:]))
			record.V[j] = math.Float32frombits(binary.LittleEndian.Uint32(wire[1040+j*4:]))
		}
		for j := range record.UVPoints {
			record.UVPoints[j] = coverageVec3(wire[1048+j*12:])
		}
		record.Origin = coverageVec3(wire[1096:])
		record.BasisNormal = coverageVec3(wire[1108:])
	}
	return records, nil
}

func EncodeWaterOverlays(records []lumps.WaterOverlay, context format.Context) ([]byte, error) {
	if err := validateCoverageContext(context, "water overlay"); err != nil {
		return nil, err
	}
	size, err := format.FixedRecordBytes(len(records), 1120)
	if err != nil {
		return nil, fmt.Errorf("water overlay: %w", err)
	}
	for i := range records {
		if count := records[i].FaceCount(); count > 256 {
			return nil, fmt.Errorf("water overlay %d: face count %d exceeds 256", i, count)
		}
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*1120 : (i+1)*1120]
		record := &records[i]
		binary.LittleEndian.PutUint32(wire, uint32(record.ID))
		binary.LittleEndian.PutUint16(wire[4:], uint16(record.TexInfo))
		binary.LittleEndian.PutUint16(wire[6:], record.FaceCountAndRenderOrder)
		for j := range record.Faces {
			binary.LittleEndian.PutUint32(wire[8+j*4:], uint32(record.Faces[j]))
		}
		for j := 0; j < 2; j++ {
			binary.LittleEndian.PutUint32(wire[1032+j*4:], math.Float32bits(record.U[j]))
			binary.LittleEndian.PutUint32(wire[1040+j*4:], math.Float32bits(record.V[j]))
		}
		for j := range record.UVPoints {
			putCoverageVec3(wire[1048+j*12:], record.UVPoints[j])
		}
		putCoverageVec3(wire[1096:], record.Origin)
		putCoverageVec3(wire[1108:], record.BasisNormal)
	}
	return raw, nil
}

func DecodeOverlaySystemLevels(raw []byte, context format.Context) ([]lumps.OverlaySystemLevel, error) {
	if err := validateCoverageModernContext(context, "overlay system level", true); err != nil {
		return nil, err
	}
	count, err := format.FixedRecordCount(len(raw), 4)
	if err != nil {
		return nil, fmt.Errorf("overlay system level: %w", err)
	}
	records := make([]lumps.OverlaySystemLevel, count)
	for i := range records {
		wire := raw[i*4 : (i+1)*4]
		records[i] = lumps.OverlaySystemLevel{MinCPULevel: wire[0], MaxCPULevel: wire[1], MinGPULevel: wire[2], MaxGPULevel: wire[3]}
	}
	return records, nil
}

func EncodeOverlaySystemLevels(records []lumps.OverlaySystemLevel, context format.Context) ([]byte, error) {
	if err := validateCoverageModernContext(context, "overlay system level", true); err != nil {
		return nil, err
	}
	size, err := format.FixedRecordBytes(len(records), 4)
	if err != nil {
		return nil, fmt.Errorf("overlay system level: %w", err)
	}
	raw := make([]byte, size)
	for i := range records {
		record := &records[i]
		raw[i*4], raw[i*4+1], raw[i*4+2], raw[i*4+3] = record.MinCPULevel, record.MaxCPULevel, record.MinGPULevel, record.MaxGPULevel
	}
	return raw, nil
}

// SDK's unused ID63 and L4D2's UNUSED_63 are not modern multiblend aliases.
func DecodeDisplacementMultiblends(raw []byte, context format.Context) ([]lumps.DisplacementMultiblend, error) {
	if err := validateCoverageModernContext(context, "displacement multiblend", false); err != nil {
		return nil, err
	}
	count, err := format.FixedRecordCount(len(raw), 80)
	if err != nil {
		return nil, fmt.Errorf("displacement multiblend: %w", err)
	}
	records := make([]lumps.DisplacementMultiblend, count)
	for i := range records {
		wire := raw[i*80 : (i+1)*80]
		record := &records[i]
		for j := 0; j < 4; j++ {
			record.Weights[j] = math.Float32frombits(binary.LittleEndian.Uint32(wire[j*4:]))
			record.Alpha[j] = math.Float32frombits(binary.LittleEndian.Uint32(wire[16+j*4:]))
			record.Colors[j] = coverageVec3(wire[32+j*12:])
		}
	}
	return records, nil
}

func EncodeDisplacementMultiblends(records []lumps.DisplacementMultiblend, context format.Context) ([]byte, error) {
	if err := validateCoverageModernContext(context, "displacement multiblend", false); err != nil {
		return nil, err
	}
	size, err := format.FixedRecordBytes(len(records), 80)
	if err != nil {
		return nil, fmt.Errorf("displacement multiblend: %w", err)
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*80 : (i+1)*80]
		record := &records[i]
		for j := 0; j < 4; j++ {
			binary.LittleEndian.PutUint32(wire[j*4:], math.Float32bits(record.Weights[j]))
			binary.LittleEndian.PutUint32(wire[16+j*4:], math.Float32bits(record.Alpha[j]))
			putCoverageVec3(wire[32+j*12:], record.Colors[j])
		}
	}
	return raw, nil
}

// A self-contained span must be nonnegative and its exclusive end must fit
// the signed index domain. Actual sibling bounds are an attached File concern.
func validateCoverageSpan(start, count int32, name string, index int) error {
	if start < 0 || count < 0 || int64(start)+int64(count) > math.MaxInt32 {
		return fmt.Errorf("%s %d: invalid span start=%d count=%d", name, index, start, count)
	}
	return nil
}

func DecodePropCollisions(raw []byte, context format.Context) ([]lumps.PropCollision, error) {
	if err := validatePropCoverageContext(context, "prop collision"); err != nil {
		return nil, err
	}
	count, err := format.FixedRecordCount(len(raw), 8)
	if err != nil {
		return nil, fmt.Errorf("prop collision: %w", err)
	}
	for i := 0; i < count; i++ {
		wire := raw[i*8 : (i+1)*8]
		if err := validateCoverageSpan(int32(binary.LittleEndian.Uint32(wire[4:])), int32(binary.LittleEndian.Uint32(wire)), "prop collision", i); err != nil {
			return nil, err
		}
	}
	records := make([]lumps.PropCollision, count)
	for i := range records {
		wire := raw[i*8 : (i+1)*8]
		records[i] = lumps.PropCollision{HullCount: int32(binary.LittleEndian.Uint32(wire)), HullStart: int32(binary.LittleEndian.Uint32(wire[4:]))}
	}
	return records, nil
}

func EncodePropCollisions(records []lumps.PropCollision, context format.Context) ([]byte, error) {
	if err := validatePropCoverageContext(context, "prop collision"); err != nil {
		return nil, err
	}
	size, err := format.FixedRecordBytes(len(records), 8)
	if err != nil {
		return nil, fmt.Errorf("prop collision: %w", err)
	}
	for i := range records {
		if err := validateCoverageSpan(records[i].HullStart, records[i].HullCount, "prop collision", i); err != nil {
			return nil, err
		}
	}
	raw := make([]byte, size)
	for i := range records {
		binary.LittleEndian.PutUint32(raw[i*8:], uint32(records[i].HullCount))
		binary.LittleEndian.PutUint32(raw[i*8+4:], uint32(records[i].HullStart))
	}
	return raw, nil
}

func DecodePropHulls(raw []byte, context format.Context) ([]lumps.PropHull, error) {
	if err := validatePropCoverageContext(context, "prop hull"); err != nil {
		return nil, err
	}
	count, err := format.FixedRecordCount(len(raw), 16)
	if err != nil {
		return nil, fmt.Errorf("prop hull: %w", err)
	}
	for i := 0; i < count; i++ {
		wire := raw[i*16 : (i+1)*16]
		if err := validateCoverageSpan(int32(binary.LittleEndian.Uint32(wire[4:])), int32(binary.LittleEndian.Uint32(wire)), "prop hull", i); err != nil {
			return nil, err
		}
	}
	records := make([]lumps.PropHull, count)
	for i := range records {
		wire := raw[i*16 : (i+1)*16]
		records[i] = lumps.PropHull{VertexCount: int32(binary.LittleEndian.Uint32(wire)), VertexStart: int32(binary.LittleEndian.Uint32(wire[4:])), SurfaceProperty: int32(binary.LittleEndian.Uint32(wire[8:])), Contents: binary.LittleEndian.Uint32(wire[12:])}
	}
	return records, nil
}

func EncodePropHulls(records []lumps.PropHull, context format.Context) ([]byte, error) {
	if err := validatePropCoverageContext(context, "prop hull"); err != nil {
		return nil, err
	}
	size, err := format.FixedRecordBytes(len(records), 16)
	if err != nil {
		return nil, fmt.Errorf("prop hull: %w", err)
	}
	for i := range records {
		if err := validateCoverageSpan(records[i].VertexStart, records[i].VertexCount, "prop hull", i); err != nil {
			return nil, err
		}
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*16 : (i+1)*16]
		record := &records[i]
		binary.LittleEndian.PutUint32(wire, uint32(record.VertexCount))
		binary.LittleEndian.PutUint32(wire[4:], uint32(record.VertexStart))
		binary.LittleEndian.PutUint32(wire[8:], uint32(record.SurfaceProperty))
		binary.LittleEndian.PutUint32(wire[12:], record.Contents)
	}
	return raw, nil
}

func DecodePropHullVertices(raw []byte, context format.Context) ([]mgl32.Vec3, error) {
	if err := validatePropCoverageContext(context, "prop hull vertex"); err != nil {
		return nil, err
	}
	count, err := format.FixedRecordCount(len(raw), 12)
	if err != nil {
		return nil, fmt.Errorf("prop hull vertex: %w", err)
	}
	records := make([]mgl32.Vec3, count)
	for i := range records {
		records[i] = coverageVec3(raw[i*12:])
	}
	return records, nil
}

func EncodePropHullVertices(records []mgl32.Vec3, context format.Context) ([]byte, error) {
	if err := validatePropCoverageContext(context, "prop hull vertex"); err != nil {
		return nil, err
	}
	size, err := format.FixedRecordBytes(len(records), 12)
	if err != nil {
		return nil, fmt.Errorf("prop hull vertex: %w", err)
	}
	raw := make([]byte, size)
	for i := range records {
		putCoverageVec3(raw[i*12:], records[i])
	}
	return raw, nil
}

func DecodePropHullTriangles(raw []byte, context format.Context) ([]lumps.PropHullTriangles, error) {
	if err := validatePropCoverageContext(context, "prop hull triangles"); err != nil {
		return nil, err
	}
	count, err := format.FixedRecordCount(len(raw), 8)
	if err != nil {
		return nil, fmt.Errorf("prop hull triangles: %w", err)
	}
	for i := 0; i < count; i++ {
		wire := raw[i*8 : (i+1)*8]
		if err := validateCoverageSpan(int32(binary.LittleEndian.Uint32(wire)), int32(binary.LittleEndian.Uint32(wire[4:])), "prop hull triangles", i); err != nil {
			return nil, err
		}
	}
	records := make([]lumps.PropHullTriangles, count)
	for i := range records {
		wire := raw[i*8 : (i+1)*8]
		records[i] = lumps.PropHullTriangles{IndexStart: int32(binary.LittleEndian.Uint32(wire)), IndexCount: int32(binary.LittleEndian.Uint32(wire[4:]))}
	}
	return records, nil
}

func EncodePropHullTriangles(records []lumps.PropHullTriangles, context format.Context) ([]byte, error) {
	if err := validatePropCoverageContext(context, "prop hull triangles"); err != nil {
		return nil, err
	}
	size, err := format.FixedRecordBytes(len(records), 8)
	if err != nil {
		return nil, fmt.Errorf("prop hull triangles: %w", err)
	}
	for i := range records {
		if err := validateCoverageSpan(records[i].IndexStart, records[i].IndexCount, "prop hull triangles", i); err != nil {
			return nil, err
		}
	}
	raw := make([]byte, size)
	for i := range records {
		binary.LittleEndian.PutUint32(raw[i*8:], uint32(records[i].IndexStart))
		binary.LittleEndian.PutUint32(raw[i*8+4:], uint32(records[i].IndexCount))
	}
	return raw, nil
}

// These temporary registry adapters select dialects from Metadata.Context,
// never from payload size. Failed decodes leave existing records untouched.
type WaterOverlays struct {
	Metadata
	data []lumps.WaterOverlay
}

func (lump *WaterOverlays) Decode(raw []byte) error {
	data, err := DecodeWaterOverlays(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}
func (lump *WaterOverlays) Encode() ([]byte, error) {
	return EncodeWaterOverlays(lump.data, lump.Context())
}
func (lump *WaterOverlays) GetData() []lumps.WaterOverlay { return lump.data }

type OverlaySystemLevels struct {
	Metadata
	data []lumps.OverlaySystemLevel
}

func (lump *OverlaySystemLevels) Decode(raw []byte) error {
	data, err := DecodeOverlaySystemLevels(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}
func (lump *OverlaySystemLevels) Encode() ([]byte, error) {
	return EncodeOverlaySystemLevels(lump.data, lump.Context())
}
func (lump *OverlaySystemLevels) GetData() []lumps.OverlaySystemLevel { return lump.data }

type DisplacementMultiblends struct {
	Metadata
	data []lumps.DisplacementMultiblend
}

func (lump *DisplacementMultiblends) Decode(raw []byte) error {
	data, err := DecodeDisplacementMultiblends(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}
func (lump *DisplacementMultiblends) Encode() ([]byte, error) {
	return EncodeDisplacementMultiblends(lump.data, lump.Context())
}
func (lump *DisplacementMultiblends) GetData() []lumps.DisplacementMultiblend { return lump.data }

type PropCollisions struct {
	Metadata
	data []lumps.PropCollision
}

func (lump *PropCollisions) Decode(raw []byte) error {
	data, err := DecodePropCollisions(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}
func (lump *PropCollisions) Encode() ([]byte, error) {
	return EncodePropCollisions(lump.data, lump.Context())
}
func (lump *PropCollisions) GetData() []lumps.PropCollision { return lump.data }

type PropHulls struct {
	Metadata
	data []lumps.PropHull
}

func (lump *PropHulls) Decode(raw []byte) error {
	data, err := DecodePropHulls(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}
func (lump *PropHulls) Encode() ([]byte, error)   { return EncodePropHulls(lump.data, lump.Context()) }
func (lump *PropHulls) GetData() []lumps.PropHull { return lump.data }

type PropHullVertices struct {
	Metadata
	data []mgl32.Vec3
}

func (lump *PropHullVertices) Decode(raw []byte) error {
	data, err := DecodePropHullVertices(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}
func (lump *PropHullVertices) Encode() ([]byte, error) {
	return EncodePropHullVertices(lump.data, lump.Context())
}
func (lump *PropHullVertices) GetData() []mgl32.Vec3 { return lump.data }

type PropHullTriangles struct {
	Metadata
	data []lumps.PropHullTriangles
}

func (lump *PropHullTriangles) Decode(raw []byte) error {
	data, err := DecodePropHullTriangles(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}
func (lump *PropHullTriangles) Encode() ([]byte, error) {
	return EncodePropHullTriangles(lump.data, lump.Context())
}
func (lump *PropHullTriangles) GetData() []lumps.PropHullTriangles { return lump.data }

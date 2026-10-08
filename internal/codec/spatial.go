package codec

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

// ErrUnsupportedPrimitiveType distinguishes an unverified primitive dialect
// from malformed supported records and profile/revision mismatches.
var ErrUnsupportedPrimitiveType = errors.New("unsupported primitive type")

func validateSpatialContext(context format.Context, name string) error {
	if err := validatePCContext(context); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	if context.LumpVersion != 0 {
		return fmt.Errorf("%s: unsupported revision %d for %s", name, context.LumpVersion, context.Profile)
	}
	return nil
}

// DecodeOverlays consumes complete 352-byte records, preserving all float bits
// and inactive face slots. Sibling references are validated by the attached File.
func DecodeOverlays(raw []byte, context format.Context) ([]lumps.Overlay, error) {
	if err := validateSpatialContext(context, "overlay"); err != nil {
		return nil, err
	}
	count, err := format.FixedRecordCount(len(raw), 352)
	if err != nil {
		return nil, fmt.Errorf("overlay: %w", err)
	}
	for i := 0; i < count; i++ {
		if faces := binary.LittleEndian.Uint16(raw[i*352+6:]) & 0x3fff; faces > 64 {
			return nil, fmt.Errorf("overlay %d: face count %d exceeds 64", i, faces)
		}
	}
	records := make([]lumps.Overlay, count)
	for i := range records {
		wire := raw[i*352 : (i+1)*352]
		record := &records[i]
		record.ID = int32(binary.LittleEndian.Uint32(wire))
		record.TexInfo = int16(binary.LittleEndian.Uint16(wire[4:]))
		record.FaceCountAndRenderOrder = binary.LittleEndian.Uint16(wire[6:])
		for j := range record.Faces {
			record.Faces[j] = int32(binary.LittleEndian.Uint32(wire[8+j*4:]))
		}
		for j := 0; j < 2; j++ {
			record.U[j] = math.Float32frombits(binary.LittleEndian.Uint32(wire[264+j*4:]))
			record.V[j] = math.Float32frombits(binary.LittleEndian.Uint32(wire[272+j*4:]))
		}
		for j := range record.UVPoints {
			for k := 0; k < 3; k++ {
				record.UVPoints[j][k] = math.Float32frombits(binary.LittleEndian.Uint32(wire[280+j*12+k*4:]))
			}
		}
		for j := 0; j < 3; j++ {
			record.Origin[j] = math.Float32frombits(binary.LittleEndian.Uint32(wire[328+j*4:]))
			record.BasisNormal[j] = math.Float32frombits(binary.LittleEndian.Uint32(wire[340+j*4:]))
		}
	}
	return records, nil
}

func EncodeOverlays(records []lumps.Overlay, context format.Context) ([]byte, error) {
	if err := validateSpatialContext(context, "overlay"); err != nil {
		return nil, err
	}
	size, err := format.FixedRecordBytes(len(records), 352)
	if err != nil {
		return nil, fmt.Errorf("overlay: %w", err)
	}
	for i := range records {
		if count := records[i].FaceCount(); count > 64 {
			return nil, fmt.Errorf("overlay %d: face count %d exceeds 64", i, count)
		}
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*352 : (i+1)*352]
		record := &records[i]
		binary.LittleEndian.PutUint32(wire, uint32(record.ID))
		binary.LittleEndian.PutUint16(wire[4:], uint16(record.TexInfo))
		binary.LittleEndian.PutUint16(wire[6:], record.FaceCountAndRenderOrder)
		for j := range record.Faces {
			binary.LittleEndian.PutUint32(wire[8+j*4:], uint32(record.Faces[j]))
		}
		for j := 0; j < 2; j++ {
			binary.LittleEndian.PutUint32(wire[264+j*4:], math.Float32bits(record.U[j]))
			binary.LittleEndian.PutUint32(wire[272+j*4:], math.Float32bits(record.V[j]))
		}
		for j := range record.UVPoints {
			for k := 0; k < 3; k++ {
				binary.LittleEndian.PutUint32(wire[280+j*12+k*4:], math.Float32bits(record.UVPoints[j][k]))
			}
		}
		for j := 0; j < 3; j++ {
			binary.LittleEndian.PutUint32(wire[328+j*4:], math.Float32bits(record.Origin[j]))
			binary.LittleEndian.PutUint32(wire[340+j*4:], math.Float32bits(record.BasisNormal[j]))
		}
	}
	return raw, nil
}

func DecodeCubemaps(raw []byte, context format.Context) ([]lumps.Cubemap, error) {
	if err := validateSpatialContext(context, "cubemap"); err != nil {
		return nil, err
	}
	count, err := format.FixedRecordCount(len(raw), 16)
	if err != nil {
		return nil, fmt.Errorf("cubemap: %w", err)
	}
	records := make([]lumps.Cubemap, count)
	for i := range records {
		wire := raw[i*16 : (i+1)*16]
		for j := 0; j < 3; j++ {
			records[i].Origin[j] = int32(binary.LittleEndian.Uint32(wire[j*4:]))
		}
		records[i].Size = wire[12]
		copy(records[i].Padding[:], wire[13:16])
	}
	return records, nil
}

func EncodeCubemaps(records []lumps.Cubemap, context format.Context) ([]byte, error) {
	if err := validateSpatialContext(context, "cubemap"); err != nil {
		return nil, err
	}
	size, err := format.FixedRecordBytes(len(records), 16)
	if err != nil {
		return nil, fmt.Errorf("cubemap: %w", err)
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*16 : (i+1)*16]
		for j := 0; j < 3; j++ {
			binary.LittleEndian.PutUint32(wire[j*4:], uint32(records[i].Origin[j]))
		}
		wire[12] = records[i].Size
		copy(wire[13:16], records[i].Padding[:])
	}
	return raw, nil
}

func DecodePrimitives(raw []byte, context format.Context) ([]lumps.Primitive, error) {
	if err := validateSpatialContext(context, "primitive"); err != nil {
		return nil, err
	}
	count, err := format.FixedRecordCount(len(raw), 10)
	if err != nil {
		return nil, fmt.Errorf("primitive: %w", err)
	}
	for i := 0; i < count; i++ {
		if kind := raw[i*10]; kind > 1 {
			return nil, fmt.Errorf("primitive %d: %w %d", i, ErrUnsupportedPrimitiveType, kind)
		}
	}
	records := make([]lumps.Primitive, count)
	for i := range records {
		wire := raw[i*10 : (i+1)*10]
		records[i] = lumps.Primitive{
			Type: wire[0], Padding: wire[1],
			FirstIndex: binary.LittleEndian.Uint16(wire[2:]),
			IndexCount: binary.LittleEndian.Uint16(wire[4:]),
			FirstVert:  binary.LittleEndian.Uint16(wire[6:]),
			VertCount:  binary.LittleEndian.Uint16(wire[8:]),
		}
	}
	return records, nil
}

func EncodePrimitives(records []lumps.Primitive, context format.Context) ([]byte, error) {
	if err := validateSpatialContext(context, "primitive"); err != nil {
		return nil, err
	}
	size, err := format.FixedRecordBytes(len(records), 10)
	if err != nil {
		return nil, fmt.Errorf("primitive: %w", err)
	}
	for i := range records {
		if records[i].Type > 1 {
			return nil, fmt.Errorf("primitive %d: %w %d", i, ErrUnsupportedPrimitiveType, records[i].Type)
		}
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*10 : (i+1)*10]
		record := &records[i]
		wire[0], wire[1] = record.Type, record.Padding
		binary.LittleEndian.PutUint16(wire[2:], record.FirstIndex)
		binary.LittleEndian.PutUint16(wire[4:], record.IndexCount)
		binary.LittleEndian.PutUint16(wire[6:], record.FirstVert)
		binary.LittleEndian.PutUint16(wire[8:], record.VertCount)
	}
	return raw, nil
}

// Overlays, Cubemaps and Primitives are temporary registry adapters. Their
// metadata context, not payload size, selects the verified wire dialect.
type Overlays struct {
	Metadata
	data []lumps.Overlay
}

func (lump *Overlays) Decode(raw []byte) error {
	data, err := DecodeOverlays(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}
func (lump *Overlays) Encode() ([]byte, error)  { return EncodeOverlays(lump.data, lump.Context()) }
func (lump *Overlays) GetData() []lumps.Overlay { return lump.data }

type Cubemaps struct {
	Metadata
	data []lumps.Cubemap
}

func (lump *Cubemaps) Decode(raw []byte) error {
	data, err := DecodeCubemaps(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}
func (lump *Cubemaps) Encode() ([]byte, error)  { return EncodeCubemaps(lump.data, lump.Context()) }
func (lump *Cubemaps) GetData() []lumps.Cubemap { return lump.data }

type Primitives struct {
	Metadata
	data []lumps.Primitive
}

func (lump *Primitives) Decode(raw []byte) error {
	data, err := DecodePrimitives(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}
func (lump *Primitives) Encode() ([]byte, error) {
	return EncodePrimitives(lump.data, lump.Context())
}
func (lump *Primitives) GetData() []lumps.Primitive { return lump.data }

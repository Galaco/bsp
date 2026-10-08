package codec

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

func validateFaceContext(context format.Context) error {
	if err := validatePCContext(context); err != nil {
		return err
	}
	if context.LumpVersion != 0 && context.LumpVersion != 1 {
		return fmt.Errorf("unsupported face revision %d", context.LumpVersion)
	}
	return nil
}

// DecodeFaces shares the invariant 56-byte layout for ordinary/original/HDR
// faces. NumPrims remains the exact packed count/shadow word.
func DecodeFaces(raw []byte, context format.Context) ([]lumps.Face, error) {
	if err := validateFaceContext(context); err != nil {
		return nil, err
	}
	count, err := format.FixedRecordCount(len(raw), 56)
	if err != nil {
		return nil, err
	}
	records := make([]lumps.Face, count)
	for i := range records {
		wire := raw[i*56 : (i+1)*56]
		record := &records[i]
		record.Planenum = binary.LittleEndian.Uint16(wire)
		record.Side, record.OnNode = wire[2], wire[3]
		record.FirstEdge = int32(binary.LittleEndian.Uint32(wire[4:]))
		record.NumEdges = int16(binary.LittleEndian.Uint16(wire[8:]))
		record.TexInfo = int16(binary.LittleEndian.Uint16(wire[10:]))
		record.DispInfo = int16(binary.LittleEndian.Uint16(wire[12:]))
		record.SurfaceFogVolumeID = int16(binary.LittleEndian.Uint16(wire[14:]))
		copy(record.Styles[:], wire[16:20])
		record.Lightofs = int32(binary.LittleEndian.Uint32(wire[20:]))
		record.Area = math.Float32frombits(binary.LittleEndian.Uint32(wire[24:]))
		for j := 0; j < 2; j++ {
			record.LightmapTextureMinsInLuxels[j] = int32(binary.LittleEndian.Uint32(wire[28+j*4:]))
			record.LightmapTextureSizeInLuxels[j] = int32(binary.LittleEndian.Uint32(wire[36+j*4:]))
		}
		record.OrigFace = int32(binary.LittleEndian.Uint32(wire[44:]))
		record.NumPrims = binary.LittleEndian.Uint16(wire[48:])
		record.FirstPrimID = binary.LittleEndian.Uint16(wire[50:])
		record.SmoothingGroups = binary.LittleEndian.Uint32(wire[52:])
	}
	return records, nil
}

func EncodeFaces(records []lumps.Face, context format.Context) ([]byte, error) {
	if err := validateFaceContext(context); err != nil {
		return nil, err
	}
	size, err := format.FixedRecordBytes(len(records), 56)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*56 : (i+1)*56]
		record := &records[i]
		binary.LittleEndian.PutUint16(wire, record.Planenum)
		wire[2], wire[3] = record.Side, record.OnNode
		binary.LittleEndian.PutUint32(wire[4:], uint32(record.FirstEdge))
		binary.LittleEndian.PutUint16(wire[8:], uint16(record.NumEdges))
		binary.LittleEndian.PutUint16(wire[10:], uint16(record.TexInfo))
		binary.LittleEndian.PutUint16(wire[12:], uint16(record.DispInfo))
		binary.LittleEndian.PutUint16(wire[14:], uint16(record.SurfaceFogVolumeID))
		copy(wire[16:20], record.Styles[:])
		binary.LittleEndian.PutUint32(wire[20:], uint32(record.Lightofs))
		binary.LittleEndian.PutUint32(wire[24:], math.Float32bits(record.Area))
		for j := 0; j < 2; j++ {
			binary.LittleEndian.PutUint32(wire[28+j*4:], uint32(record.LightmapTextureMinsInLuxels[j]))
			binary.LittleEndian.PutUint32(wire[36+j*4:], uint32(record.LightmapTextureSizeInLuxels[j]))
		}
		binary.LittleEndian.PutUint32(wire[44:], uint32(record.OrigFace))
		binary.LittleEndian.PutUint16(wire[48:], record.NumPrims)
		binary.LittleEndian.PutUint16(wire[50:], record.FirstPrimID)
		binary.LittleEndian.PutUint32(wire[52:], record.SmoothingGroups)
	}
	return raw, nil
}

// Faces is the shared private registry adapter for complete owned face records.
type Faces struct {
	Metadata
	data []lumps.Face
}

func (lump *Faces) Decode(raw []byte) error {
	data, err := DecodeFaces(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}
func (lump *Faces) Encode() ([]byte, error) { return EncodeFaces(lump.data, lump.Context()) }
func (lump *Faces) GetData() []lumps.Face   { return lump.data }

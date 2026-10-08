package bsp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/galaco/bsp/lumps"
	"github.com/galaco/bsp/primitives/plane"
)

// Bsp is the root .bsp filetype container.
// Consists of a 1036byte header and 64 lump blocks.
type Bsp struct {
	header Header
	lumps  [LumpCount]payload
}

// Header is the Bsp header. Contains format and lump layout data.
// Do not trust lump information between import and export
type Header struct {
	Id       int32
	Version  int32
	Lumps    [64]HeaderLump
	Revision int32
}

// HeaderLump contains layout information for a given lump, stored in the Header.
type HeaderLump struct {
	Offset  int32
	Length  int32
	Version int32
	Id      [4]byte
}

// Header returns metadata by value. Offsets describe the input, not a future export.
func (bsp *Bsp) Header() Header { return bsp.header }

// Descriptor returns current stored length and original location metadata.
func (bsp *Bsp) Descriptor(index LumpId) (HeaderLump, error) {
	if err := validateID(index); err != nil {
		return HeaderLump{}, err
	}
	return bsp.header.Lumps[index], nil
}

// Lump returns an independent decoded snapshot. Changes require an explicit setter.
func (bsp *Bsp) Lump(index LumpId) (lumps.ILump, error) {
	if err := validateID(index); err != nil {
		return nil, err
	}
	decoder, err := getReferenceLumpByIndex(int(index), bsp.header.Version)
	if err != nil {
		return nil, err
	}
	raw := bsp.lumps[index].raw
	switch decoder.(type) {
	case *lumps.Generic, *lumps.Unimplemented, *lumps.Game, *lumps.Visibility,
		*lumps.PhysDisp, *lumps.DispLightmapSamplePosition:
		raw = cloneBytes(raw)
	}
	if game, ok := decoder.(*lumps.Game); ok {
		game.UpdateInternalOffsets(bsp.header.Lumps[index].Offset)
	}
	if err := decoder.Unmarshall(raw); err != nil {
		return nil, &FormatError{fmt.Sprintf("decode lump %d", index), err}
	}
	return decoder, nil
}

// RawLump returns an independent snapshot of stored bytes.
func (bsp *Bsp) RawLump(index LumpId) ([]byte, error) {
	if err := validateID(index); err != nil {
		return nil, err
	}
	return cloneBytes(bsp.lumps[index].raw), nil
}

// SetRawLump atomically owns bytes and metadata. Length is derived, not caller supplied.
// Offset records the encoded input base, particularly for game-directory offsets.
func (bsp *Bsp) SetRawLump(index LumpId, descriptor HeaderLump, raw []byte) error {
	if err := validateID(index); err != nil {
		return err
	}
	if descriptor.Offset < 0 || int64(len(raw)) > math.MaxInt32 {
		return fmt.Errorf("%w: invalid replacement range", ErrMalformed)
	}
	bsp.storeOwned(index, descriptor, cloneBytes(raw))
	return nil
}

func (bsp *Bsp) storeOwned(index LumpId, descriptor HeaderLump, owned []byte) {
	descriptor.Length = int32(len(owned))
	bsp.lumps[index] = payload{raw: owned}
	bsp.header.Lumps[index] = descriptor
}

// Entities returns an independent entity-text snapshot.
func (bsp *Bsp) Entities() (string, error) {
	data, err := bsp.Lump(LumpEntities)
	if err != nil {
		return "", err
	}
	entities, ok := data.(*lumps.EntData)
	if !ok {
		return "", ErrUnsupportedLump
	}
	return entities.GetData(), nil
}

// SetEntities commits exact entity bytes.
func (bsp *Bsp) SetEntities(text string) error {
	if int64(len(text)) > math.MaxInt32 {
		return fmt.Errorf("%w: entity payload too large", ErrMalformed)
	}
	descriptor := bsp.header.Lumps[LumpEntities]
	descriptor.Id = [4]byte{}
	bsp.storeOwned(LumpEntities, descriptor, []byte(text))
	return nil
}

// Planes returns an independent plane snapshot.
func (bsp *Bsp) Planes() ([]plane.Plane, error) {
	data, err := bsp.Lump(LumpPlanes)
	if err != nil {
		return nil, err
	}
	planes, ok := data.(*lumps.Planes)
	if !ok {
		return nil, ErrUnsupportedLump
	}
	return planes.GetData(), nil
}

// SetPlanes validates and commits plane records without retaining caller slices.
func (bsp *Bsp) SetPlanes(planes []plane.Plane) error {
	if len(planes) > 65536 {
		return fmt.Errorf("%w: too many planes", ErrMalformed)
	}
	for _, p := range planes {
		if p.AxisType < 0 || p.AxisType > 5 {
			return fmt.Errorf("%w: plane axis %d", ErrMalformed, p.AxisType)
		}
	}
	var out bytes.Buffer
	if err := binary.Write(&out, binary.LittleEndian, planes); err != nil {
		return err
	}
	descriptor := bsp.header.Lumps[LumpPlanes]
	descriptor.Id = [4]byte{}
	bsp.storeOwned(LumpPlanes, descriptor, out.Bytes())
	return nil
}

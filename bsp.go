package bsp

import (
	"errors"
	"fmt"
	"math"

	"github.com/galaco/bsp/internal/codec"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

// File is an owned BSP container with a header and 64 lump payloads.
// Concurrent access involving edits requires external synchronization.
type File struct {
	header  Header
	lumps   [LumpCount]payload
	profile Profile
	layout  format.HeaderLayout
	limits  codec.Limits
}

// Header contains format metadata and current stored lump lengths.
// Offsets describe historical input locations, not the next export layout.
type Header struct {
	ID       int32
	Version  int32
	Lumps    [64]HeaderLump
	Revision int32
}

// HeaderLump contains layout information for a given lump, stored in the Header.
type HeaderLump struct {
	Offset  int32
	Length  int32
	Version int32
	FourCC  [4]byte // Fourth header word: compressed uncompressed-size metadata, not a lump ID.
}

// Header returns metadata by value. Offsets describe the input, not a future export.
func (bsp *File) Header() Header { return bsp.header }

// Descriptor returns current stored length and original location metadata.
func (bsp *File) Descriptor(index LumpID) (HeaderLump, error) {
	if err := validateID(index); err != nil {
		return HeaderLump{}, err
	}
	return bsp.header.Lumps[index], nil
}

// decodeLump returns an independent decoded snapshot for attached typed getters.
func (bsp *File) decodeLump(index LumpID) (codec.Decoder, error) {
	if err := validateID(index); err != nil {
		return nil, err
	}
	context := bsp.codecContext(index)
	decoder, err := codec.NewDecoder(int(index), context)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedLump, err)
	}
	raw, err := bsp.decodedLumpView(index)
	if err != nil {
		return nil, err
	}
	if err := decoder.Decode(raw); err != nil {
		if errors.Is(err, codec.ErrUnsupportedPrimitiveType) {
			return nil, fmt.Errorf("%w: %s", ErrUnsupportedLump, err)
		}
		return nil, &FormatError{fmt.Sprintf("decode lump %d", index), err}
	}
	if err := bsp.validateOrdinarySnapshot(index, decoder); err != nil {
		return nil, err
	}
	return decoder, nil
}

// RawLump returns an independent snapshot of stored bytes.
func (bsp *File) RawLump(index LumpID) ([]byte, error) {
	if err := validateID(index); err != nil {
		return nil, err
	}
	return cloneBytes(bsp.lumps[index].raw), nil
}

// SetRawLump atomically owns bytes and metadata. Length is derived, not caller supplied.
// Offset records the encoded input base, particularly for game-directory offsets.
func (bsp *File) SetRawLump(index LumpID, descriptor HeaderLump, raw []byte) error {
	if err := validateID(index); err != nil {
		return err
	}
	if descriptor.Offset < 0 || int64(len(raw)) > math.MaxInt32 {
		return fmt.Errorf("%w: invalid replacement range", ErrMalformed)
	}
	bsp.storeOwned(index, descriptor, cloneBytes(raw))
	return nil
}

func (bsp *File) storeOwned(index LumpID, descriptor HeaderLump, owned []byte) {
	descriptor.Length = int32(len(owned))
	bsp.lumps[index] = payload{raw: owned}
	bsp.header.Lumps[index] = descriptor
}

// Entities returns an independent entity-text snapshot.
func (bsp *File) Entities() (string, error) {
	data, err := bsp.decodeLump(LumpEntities)
	if err != nil {
		return "", err
	}
	entities, ok := data.(*codec.EntData)
	if !ok {
		return "", ErrUnsupportedLump
	}
	return entities.GetData(), nil
}

// SetEntities commits exact entity bytes.
func (bsp *File) SetEntities(text string) error {
	if err := bsp.checkTypedRevision(LumpEntities); err != nil {
		return err
	}
	if int64(len(text)) > math.MaxInt32 {
		return fmt.Errorf("%w: entity payload too large", ErrMalformed)
	}
	descriptor := bsp.header.Lumps[LumpEntities]
	descriptor.FourCC = [4]byte{}
	bsp.storeOwned(LumpEntities, descriptor, []byte(text))
	return nil
}

// Planes returns an independent plane snapshot.
func (bsp *File) Planes() ([]lumps.Plane, error) {
	data, err := bsp.decodeLump(LumpPlanes)
	if err != nil {
		return nil, err
	}
	planes, ok := data.(*codec.Planes)
	if !ok {
		return nil, ErrUnsupportedLump
	}
	return planes.GetData(), nil
}

// SetPlanes validates and commits plane records without retaining caller slices.
func (bsp *File) SetPlanes(planes []lumps.Plane) error {
	if err := bsp.checkTypedRevision(LumpPlanes); err != nil {
		return err
	}
	if len(planes) > 65536 {
		return fmt.Errorf("%w: too many planes", ErrMalformed)
	}
	for _, p := range planes {
		if p.AxisType < 0 || p.AxisType > 5 {
			return fmt.Errorf("%w: plane axis %d", ErrMalformed, p.AxisType)
		}
	}
	raw, err := codec.EncodePlanes(planes)
	if err != nil {
		return &FormatError{"encode planes", err}
	}
	return bsp.commitOrdinary(LumpPlanes, raw, nil)
}

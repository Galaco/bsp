package bsp

import (
	"encoding/binary"
	"fmt"
	"github.com/galaco/bsp/internal/codec"
	"github.com/galaco/bsp/internal/format"
)

// Profile identifies an engine's wire dialect. BSP version alone is insufficient.
type Profile string

const (
	ProfileSDK2013 Profile = "sdk2013"
	ProfileCSGO    Profile = "csgo"
	ProfileL4D2    Profile = "l4d2"
)

// ReadOptions selects a dialect and overridable lazy-decode safety limits.
// Zero limits select defaults; they are library policy, not format maxima.
type ReadOptions struct {
	Profile                Profile
	MaxDecodedLumpBytes    int64
	MaxLZMADictionaryBytes int64
}

const (
	DefaultMaxDecodedLumpBytes    = codec.DefaultMaxDecodedLumpBytes
	DefaultMaxLZMADictionaryBytes = codec.DefaultMaxLZMADictionaryBytes
)

func (options ReadOptions) limits() (codec.Limits, error) {
	if options.MaxDecodedLumpBytes < 0 || options.MaxLZMADictionaryBytes < 0 {
		return codec.Limits{}, fmt.Errorf("%w: negative decode limit", ErrInvalidOptions)
	}
	output, dictionary := uint64(options.MaxDecodedLumpBytes), uint64(options.MaxLZMADictionaryBytes)
	if output == 0 {
		output = DefaultMaxDecodedLumpBytes
	}
	if dictionary == 0 {
		dictionary = DefaultMaxLZMADictionaryBytes
	}
	return codec.Limits{MaxDecodedLumpBytes: output, MaxLZMADictionaryBytes: dictionary}, nil
}

func profileLayout(profile Profile, version int32) (format.HeaderLayout, error) {
	switch profile {
	case ProfileSDK2013:
		if version == 19 || version == 20 {
			return format.StandardHeader, nil
		}
	case ProfileCSGO:
		if version == 21 {
			return format.StandardHeader, nil
		}
	case ProfileL4D2:
		if version == 21 {
			return format.L4D2Header, nil
		}
	}
	return 0, fmt.Errorf("%w: profile %q, BSP version %d", ErrUnsupportedFormat, profile, version)
}

func decodeHeader(raw []byte, layout format.HeaderLayout) Header {
	var header Header
	order := binary.LittleEndian
	header.Id = int32(order.Uint32(raw))
	header.Version = int32(order.Uint32(raw[4:]))
	for i := 0; i < LumpCount; i++ {
		entry := raw[8+i*16 : 8+(i+1)*16]
		desc := &header.Lumps[i]
		if layout == format.L4D2Header {
			desc.Version = int32(order.Uint32(entry))
			desc.Offset = int32(order.Uint32(entry[4:]))
			desc.Length = int32(order.Uint32(entry[8:]))
		} else {
			desc.Offset = int32(order.Uint32(entry))
			desc.Length = int32(order.Uint32(entry[4:]))
			desc.Version = int32(order.Uint32(entry[8:]))
		}
		copy(desc.FourCC[:], entry[12:16])
	}
	header.Revision = int32(order.Uint32(raw[1032:]))
	return header
}

func encodeHeader(raw []byte, header *Header, layout format.HeaderLayout) {
	order := binary.LittleEndian
	order.PutUint32(raw[:], uint32(header.Id))
	order.PutUint32(raw[4:], uint32(header.Version))
	for i := range header.Lumps {
		desc := &header.Lumps[i]
		entry := raw[8+i*16 : 8+(i+1)*16]
		if layout == format.L4D2Header {
			order.PutUint32(entry, uint32(desc.Version))
			order.PutUint32(entry[4:], uint32(desc.Offset))
			order.PutUint32(entry[8:], uint32(desc.Length))
		} else {
			order.PutUint32(entry, uint32(desc.Offset))
			order.PutUint32(entry[4:], uint32(desc.Length))
			order.PutUint32(entry[8:], uint32(desc.Version))
		}
		copy(entry[12:], desc.FourCC[:])
	}
	order.PutUint32(raw[1032:], uint32(header.Revision))
}

func (bsp *Bsp) codecContext(index LumpId) format.Context {
	return format.Context{
		Profile: string(bsp.profile), BSPVersion: bsp.header.Version,
		LumpVersion: bsp.header.Lumps[index].Version,
		ByteOrder:   binary.LittleEndian, HeaderLayout: bsp.layout,
		GameOffsets:     format.AbsoluteGameOffsets,
		CompressionSize: binary.LittleEndian.Uint32(bsp.header.Lumps[index].FourCC[:]),
		LegacyAmbient:   bsp.header.Lumps[LumpLeafAmbientIndex].Length == 0,
	}
}

// Profile returns the explicit dialect used to interpret this file.
func (bsp *Bsp) Profile() Profile { return bsp.profile }

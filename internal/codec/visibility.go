package codec

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

// Visibility is the context-aware private adapter for lump 4.
type Visibility struct {
	lumps.Metadata
	data lumps.Visibility
}

func (lump *Visibility) Unmarshall(raw []byte) error {
	value, err := DecodeVisibility(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = value
	lump.SetLength(len(raw))
	return nil
}

func (lump *Visibility) Marshall() ([]byte, error) {
	return EncodeVisibility(lump.data, lump.Context())
}

func (lump *Visibility) GetData() lumps.Visibility { return lump.data }

// DecodeVisibility derives the table extent before allocating, then validates
// every referenced PVS/PAS row before publishing an owned value.
func DecodeVisibility(raw []byte, context format.Context) (lumps.Visibility, error) {
	if err := visibilityContext(context); err != nil {
		return lumps.Visibility{}, err
	}
	if len(raw) < 4 {
		return lumps.Visibility{}, fmt.Errorf("visibility: cluster count: %w", io.ErrUnexpectedEOF)
	}
	count := int32(binary.LittleEndian.Uint32(raw[:4]))
	tableEnd, err := format.VisibilitySize(int(count), 0)
	if err != nil {
		return lumps.Visibility{}, err
	}
	if tableEnd > len(raw) {
		return lumps.Visibility{}, fmt.Errorf("visibility: cluster count %d exceeds available offset table: %w", count, io.ErrUnexpectedEOF)
	}
	if _, err := format.VisibilitySize(int(count), len(raw)-tableEnd); err != nil {
		return lumps.Visibility{}, err
	}
	offsets := make([][2]int32, int(count))
	for cluster := range offsets {
		for kind := range offsets[cluster] {
			start := 4 + 8*cluster + 4*kind
			offsets[cluster][kind] = int32(binary.LittleEndian.Uint32(raw[start : start+4]))
		}
	}
	data := raw[tableEnd:]
	if err := format.ValidateVisibility(offsets, data); err != nil {
		return lumps.Visibility{}, err
	}
	owned := make([]byte, len(data))
	copy(owned, data)
	return lumps.Visibility{Offsets: offsets, Data: owned}, nil
}

// EncodeVisibility emits the derived count, unchanged signed lump-relative
// offsets, and compressed data together. Output allocation follows validation.
func EncodeVisibility(value lumps.Visibility, context format.Context) ([]byte, error) {
	if err := visibilityContext(context); err != nil {
		return nil, err
	}
	size, err := format.VisibilitySize(len(value.Offsets), len(value.Data))
	if err != nil {
		return nil, err
	}
	if err := format.ValidateVisibility(value.Offsets, value.Data); err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	binary.LittleEndian.PutUint32(raw[:4], uint32(len(value.Offsets)))
	for cluster, pair := range value.Offsets {
		for kind, offset := range pair {
			start := 4 + 8*cluster + 4*kind
			binary.LittleEndian.PutUint32(raw[start:start+4], uint32(offset))
		}
	}
	copy(raw[4+8*len(value.Offsets):], value.Data)
	return raw, nil
}

func visibilityContext(context format.Context) error {
	if err := validatePCContext(context); err != nil {
		return fmt.Errorf("visibility: %w", err)
	}
	if context.LumpVersion != 0 {
		return fmt.Errorf("visibility: unsupported revision %d", context.LumpVersion)
	}
	return nil
}

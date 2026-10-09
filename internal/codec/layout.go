package codec

import (
	"fmt"

	"github.com/galaco/bsp/internal/format"
)

// Layout describes canonical, numeric-ID-ordered lump positions. Lengths omit
// alignment padding; Size includes the header and final alignment padding.
type Layout struct {
	Offsets [format.LumpCount]int32
	Lengths [format.LumpCount]int32
	Size    int32
}

// PlanLayout computes a complete layout without allocating payload storage.
// Empty lumps have zero offsets and lengths. Every nonempty lump starts on a
// four-byte boundary, and the complete file must fit a signed 32-bit size.
// Invalid lengths or extents return an empty plan.
func PlanLayout(lengths [format.LumpCount]int64) (Layout, error) {
	const maxSize int64 = 1<<31 - 1
	var layout Layout
	cursor := int64(format.HeaderSize)
	for id, length := range lengths {
		if length < 0 || length > maxSize {
			return Layout{}, fmt.Errorf("lump %d: length %d is outside signed 32-bit range", id, length)
		}
		if length == 0 {
			continue
		}
		if cursor > maxSize || length > maxSize-cursor {
			return Layout{}, fmt.Errorf("lump %d: extent exceeds signed 32-bit file size", id)
		}
		offset := cursor
		cursor += length
		padding := (4 - cursor%4) % 4
		if padding > maxSize-cursor {
			return Layout{}, fmt.Errorf("lump %d: alignment exceeds signed 32-bit file size", id)
		}
		cursor += padding
		layout.Offsets[id] = int32(offset)
		layout.Lengths[id] = int32(length)
	}
	layout.Size = int32(cursor)
	return layout, nil
}

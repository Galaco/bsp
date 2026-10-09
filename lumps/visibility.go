package lumps

import (
	"math/bits"

	"github.com/galaco/bsp/internal/format"
)

// VisibilityKind selects a potentially visible or potentially audible set.
type VisibilityKind uint8

const (
	PVS VisibilityKind = iota
	PAS
)

// Visibility is an independent snapshot of lump 4. The cluster count is
// len(Offsets), at most 65536. Offsets are signed lump-relative positions,
// including the four-byte count and eight bytes per offset pair. -1 alone
// means a missing row. Shared PVS/PAS offsets are valid.
//
// Data contains only compressed bytes after the offset table. Changing the
// table length requires explicitly adjusting offsets; they are not rebased.
type Visibility struct {
	Offsets [][2]int32
	Data    []byte
}

// Row returns exactly ceil(len(Offsets)/8) bytes for the selected row. Missing
// rows return nil, false, nil, rather than inventing visibility. Errors identify
// invalid indices, selectors, offsets, or malformed bounded RLE streams.
func (value Visibility) Row(cluster int, kind VisibilityKind) ([]byte, bool, error) {
	return format.VisibilityRow(value.Offsets, value.Data, cluster, int(kind))
}

// VisibleClusters returns ascending integer cluster IDs, including ID 65535.
// Unused high bits in a non-byte-aligned row are ignored. Presence has the same
// meaning as Row; a present row with no set bits is distinct from a missing row.
func (value Visibility) VisibleClusters(cluster int, kind VisibilityKind) ([]int, bool, error) {
	row, present, err := value.Row(cluster, kind)
	if err != nil || !present {
		return nil, present, err
	}
	if remainder := len(value.Offsets) % 8; remainder != 0 {
		row[len(row)-1] &= byte(1<<uint(remainder) - 1)
	}
	count := 0
	for _, word := range row {
		count += bits.OnesCount8(word)
	}
	visible := make([]int, 0, count)
	for byteIndex, word := range row {
		for word != 0 {
			visible = append(visible, byteIndex*8+bits.TrailingZeros8(word))
			word &= word - 1
		}
	}
	return visible, true, nil
}

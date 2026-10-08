package format

import (
	"fmt"
	"io"
)

// MaxVisibilityClusters is MAX_MAP_CLUSTERS in Valve SDK2013 bspfile.h,
// revision b8cfb12c0e083a2ef5b2f9f9b50f3902fa034474.
const MaxVisibilityClusters = 65536

// VisibilitySize checks the derived table size and signed BSP payload budget.
func VisibilitySize(count, dataLength int) (int, error) {
	if count < 0 || count > MaxVisibilityClusters {
		return 0, fmt.Errorf("visibility: cluster count %d outside [0,%d]", count, MaxVisibilityClusters)
	}
	tableEnd := 4 + 8*count
	if dataLength < 0 || int64(dataLength) > int64(1<<31-1)-int64(tableEnd) {
		return 0, fmt.Errorf("visibility: payload exceeds signed 32-bit size")
	}
	return tableEnd + dataLength, nil
}

// ValidateVisibility checks both columns, walking each shared compressed row
// only once and without allocating decompressed rows. Unreferenced data is opaque.
func ValidateVisibility(offsets [][2]int32, data []byte) error {
	if _, err := VisibilitySize(len(offsets), len(data)); err != nil {
		return err
	}
	tableEnd := 4 + 8*len(offsets)
	rowLength := (len(offsets) + 7) / 8
	seen := make(map[int32]struct{})
	for cluster, pair := range offsets {
		for kind, offset := range pair {
			if offset == -1 {
				continue
			}
			if _, ok := seen[offset]; ok {
				continue
			}
			start, err := visibilityDataOffset(offset, tableEnd, len(data))
			if err != nil {
				return fmt.Errorf("visibility: cluster %d kind %d: %w", cluster, kind, err)
			}
			if err := walkVisibilityRow(data, start, rowLength, nil); err != nil {
				return fmt.Errorf("visibility: cluster %d kind %d offset %d: %w", cluster, kind, offset, err)
			}
			seen[offset] = struct{}{}
		}
	}
	return nil
}

// VisibilityRow decodes one PVS (0) or PAS (1) row. Missing (-1) offsets do
// not produce a synthetic row. Offsets always remain relative to the lump start.
func VisibilityRow(offsets [][2]int32, data []byte, cluster, kind int) ([]byte, bool, error) {
	if _, err := VisibilitySize(len(offsets), len(data)); err != nil {
		return nil, false, err
	}
	if cluster < 0 || cluster >= len(offsets) {
		return nil, false, fmt.Errorf("visibility: cluster %d outside [0,%d)", cluster, len(offsets))
	}
	if kind < 0 || kind > 1 {
		return nil, false, fmt.Errorf("visibility: cluster %d invalid kind %d (want PVS=0 or PAS=1)", cluster, kind)
	}
	offset := offsets[cluster][kind]
	if offset == -1 {
		return nil, false, nil
	}
	start, err := visibilityDataOffset(offset, 4+8*len(offsets), len(data))
	if err != nil {
		return nil, false, fmt.Errorf("visibility: cluster %d kind %d: %w", cluster, kind, err)
	}
	row := make([]byte, (len(offsets)+7)/8)
	if err := walkVisibilityRow(data, start, len(row), row); err != nil {
		return nil, false, fmt.Errorf("visibility: cluster %d kind %d offset %d: %w", cluster, kind, offset, err)
	}
	return row, true, nil
}

func visibilityDataOffset(offset int32, tableEnd, dataLength int) (int, error) {
	if int64(offset) < int64(tableEnd) || int64(offset)-int64(tableEnd) >= int64(dataLength) {
		return 0, fmt.Errorf("offset %d must point into data [%d,%d); -1 is the only missing-row sentinel", offset, tableEnd, tableEnd+dataLength)
	}
	return int(offset) - tableEnd, nil
}

// The SDK CompressVis grammar is a nonzero literal or zero plus a positive
// byte-count run. Every iteration consumes input and advances bounded output.
func walkVisibilityRow(data []byte, start, rowLength int, output []byte) error {
	for input, written := start, 0; written < rowLength; {
		if input >= len(data) {
			return fmt.Errorf("row truncated after %d of %d bytes: %w", written, rowLength, io.ErrUnexpectedEOF)
		}
		literal := data[input]
		input++
		if literal != 0 {
			if output != nil {
				output[written] = literal
			}
			written++
			continue
		}
		if input >= len(data) {
			return fmt.Errorf("zero run missing count after %d of %d bytes: %w", written, rowLength, io.ErrUnexpectedEOF)
		}
		run := int(data[input])
		input++
		if run == 0 || run > rowLength-written {
			return fmt.Errorf("zero run %d invalid with %d row bytes remaining", run, rowLength-written)
		}
		written += run
	}
	return nil
}

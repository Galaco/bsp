package format

import "fmt"

// FixedRecordCount rejects partial records before a decoder allocates its slice.
func FixedRecordCount(length, width int) (int, error) {
	if length < 0 {
		return 0, fmt.Errorf("negative fixed-record length %d", length)
	}
	if width <= 0 {
		return 0, fmt.Errorf("invalid fixed-record width %d", width)
	}
	if length%width != 0 {
		return 0, fmt.Errorf("fixed-record length %d is not a multiple of width %d", length, width)
	}
	return length / width, nil
}

// FixedRecordBytes checks the BSP signed-32-bit and host address-space budgets
// before multiplication or allocation.
func FixedRecordBytes(count, width int) (int, error) {
	if count < 0 {
		return 0, fmt.Errorf("negative fixed-record count %d", count)
	}
	if width <= 0 {
		return 0, fmt.Errorf("invalid fixed-record width %d", width)
	}
	limit := uint64(1<<31 - 1)
	if addressLimit := uint64(^uint(0) >> 1); addressLimit < limit {
		limit = addressLimit
	}
	if uint64(count) > limit/uint64(width) {
		return 0, fmt.Errorf("fixed-record count %d at width %d exceeds byte budget %d", count, width, limit)
	}
	return count * width, nil
}

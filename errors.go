package bsp

import (
	"errors"
	"fmt"
)

var (
	ErrMalformed         = errors.New("malformed BSP data")
	ErrInvalidLumpID     = errors.New("invalid lump ID")
	ErrUnsupportedFormat = errors.New("unsupported BSP format")
	ErrUnsupportedLump   = errors.New("unsupported lump layout")
	ErrLumpNotFound      = errors.New("lump not found")
)

// FormatError identifies malformed data while retaining its underlying error.
type FormatError struct {
	Context string
	Err     error
}

func (e *FormatError) Error() string        { return fmt.Sprintf("%s: %v", e.Context, e.Err) }
func (e *FormatError) Unwrap() error        { return e.Err }
func (e *FormatError) Is(target error) bool { return target == ErrMalformed }

const (
	HeaderSize       = 1036
	LumpCount        = 64
	fileIdent  int32 = 0x50534256
)

func validateID(id LumpId) error {
	if id < 0 || id >= LumpCount {
		return fmt.Errorf("%w: %d", ErrInvalidLumpID, id)
	}
	return nil
}

package bsp

import (
	"errors"
	"fmt"
	"github.com/galaco/bsp/internal/codec"
	"github.com/galaco/bsp/internal/format"
)

var (
	ErrMalformed                = errors.New("malformed BSP data")
	ErrInvalidLumpID            = errors.New("invalid lump ID")
	ErrUnsupportedFormat        = errors.New("unsupported BSP format")
	ErrUnsupportedLump          = errors.New("unsupported lump layout")
	ErrLumpNotFound             = errors.New("lump not found")
	ErrAmbiguousGameLumpID      = errors.New("ambiguous game lump ID")
	ErrInvalidGameLumpIndex     = errors.New("invalid game lump index")
	ErrInvalidDisplacementIndex = errors.New("invalid displacement index")
	ErrInvalidOptions           = errors.New("invalid read options")
	ErrLimitExceeded            = codec.ErrLimitExceeded
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
	HeaderSize       = format.HeaderSize
	LumpCount        = format.LumpCount
	fileIdent  int32 = 0x50534256
)

func validateID(id LumpId) error {
	if id < 0 || id >= LumpCount {
		return fmt.Errorf("%w: %d", ErrInvalidLumpID, id)
	}
	return nil
}

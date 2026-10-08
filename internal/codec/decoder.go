package codec

import "github.com/galaco/bsp/internal/format"

// Decoder is the internal typed wire adapter selected for an ordinary lump.
type Decoder interface {
	Decode([]byte) error
	Encode() ([]byte, error)
	SetContext(format.Context)
}

// Metadata retains import size and the format context for a private adapter.
type Metadata struct {
	length  int
	context format.Context
}

func (info *Metadata) Length() int                       { return info.length }
func (info *Metadata) SetLength(length int)              { info.length = length }
func (info *Metadata) Context() format.Context           { return info.context }
func (info *Metadata) SetContext(context format.Context) { info.context = context }

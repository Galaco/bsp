package bsp

import (
	"fmt"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/internal/versions"
	"github.com/galaco/bsp/lumps"
)

type payload struct{ raw []byte }

func cloneBytes(raw []byte) []byte { return append([]byte(nil), raw...) }

// getReferenceLumpByIndex returns an independently selected codec.
func getReferenceLumpByIndex(index int, context format.Context) (lumps.ILump, error) {
	if index < 0 || index > 63 {
		return nil, fmt.Errorf("invalid lump id: %d provided", index)
	}

	l, err := versions.New(index, context)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedLump, err)
	}
	return l, nil
}

func (bsp *Bsp) checkTypedRevision(index LumpId) error {
	if err := versions.ValidateRevision(int(index), bsp.codecContext(index)); err != nil {
		return fmt.Errorf("%w: %s", ErrUnsupportedLump, err)
	}
	return nil
}

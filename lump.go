package bsp

import (
	"fmt"
	"github.com/galaco/bsp/internal/codec"
)

type payload struct{ raw []byte }

func cloneBytes(raw []byte) []byte { return append([]byte(nil), raw...) }

func (bsp *File) checkTypedRevision(index LumpID) error {
	if err := codec.ValidateRevision(int(index), bsp.codecContext(index)); err != nil {
		return fmt.Errorf("%w: %s", ErrUnsupportedLump, err)
	}
	return nil
}

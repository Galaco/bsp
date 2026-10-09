package bsp

import (
	"fmt"

	"github.com/galaco/bsp/internal/codec"
	"github.com/galaco/bsp/lumps"
)

// Visibility returns an owned offset-table and compressed-data snapshot. An
// absent lump returns ErrLumpNotFound; a present zero-count record is distinct.
func (bsp *File) Visibility() (lumps.Visibility, error) {
	if err := bsp.checkTypedRevision(LumpVisibility); err != nil {
		return lumps.Visibility{}, err
	}
	if len(bsp.lumps[LumpVisibility].raw) == 0 {
		return lumps.Visibility{}, fmt.Errorf("%w: visibility", ErrLumpNotFound)
	}
	decoded, err := bsp.decodeLump(LumpVisibility)
	if err != nil {
		return lumps.Visibility{}, err
	}
	return decoded.(*codec.Visibility).GetData(), nil
}

// SetVisibility validates both columns and every referenced row before atomic
// replacement. Count is derived from Offsets; a zero value emits four bytes.
// Offsets stay lump-relative, so table-length edits require explicit rebasing.
func (bsp *File) SetVisibility(value lumps.Visibility) error {
	if err := bsp.checkTypedRevision(LumpVisibility); err != nil {
		return err
	}
	raw, err := codec.EncodeVisibility(value, bsp.codecContext(LumpVisibility))
	return bsp.commitOrdinary(LumpVisibility, raw, err)
}

package bsp

import (
	"fmt"
	"github.com/galaco/bsp/internal/codec"
	"github.com/galaco/bsp/internal/format"
)

// Writer binds a file explicitly; export never mutates authoritative bytes.
type Writer struct{ data *Bsp }

type exportPlan struct {
	layout       codec.Layout
	headerLayout format.HeaderLayout
	game         codec.GameRelocation
}

// prepareExport is shared by emission and current-state map checksums.
func (file *Bsp) prepareExport() (exportPlan, error) {
	if file == nil || file.header.Id != fileIdent {
		return exportPlan{}, fmt.Errorf("%w: nil file or invalid file identity", ErrMalformed)
	}
	layoutKind, err := profileLayout(file.profile, file.header.Version)
	if err != nil {
		return exportPlan{}, err
	}
	var lengths [LumpCount]int64
	for i := range lengths {
		lengths[i] = int64(len(file.lumps[i].raw))
	}
	layout, err := codec.PlanLayout(lengths)
	if err != nil {
		return exportPlan{}, &FormatError{"output layout", err}
	}
	var game codec.GameRelocation
	if lengths[LumpGame] != 0 {
		view, err := file.gameDirectoryView()
		if err != nil {
			return exportPlan{}, err
		}
		game, err = codec.PlanGameRelocation(view, layout.Offsets[LumpGame])
		if err != nil {
			return exportPlan{}, &FormatError{"output game directory", err}
		}
	}
	return exportPlan{layout: layout, headerLayout: layoutKind, game: game}, nil
}

// Write emits a canonical numeric-ID layout with aligned, lossless stored payloads.
// Historical input positions are not changed; game directory offsets are relocated.
func (w *Writer) Write() ([]byte, error) {
	if w == nil || w.data == nil {
		return nil, fmt.Errorf("%w: nil file", ErrMalformed)
	}
	file := w.data
	plan, err := file.prepareExport()
	if err != nil {
		return nil, err
	}
	layout, game := plan.layout, plan.game
	// The entire layout and game fixups are checked before the sole output allocation.
	raw := make([]byte, int(layout.Size))
	header := file.header
	for i := range header.Lumps {
		header.Lumps[i].Offset = layout.Offsets[i]
		header.Lumps[i].Length = layout.Lengths[i]
		if layout.Lengths[i] == 0 {
			continue
		}
		start := int(layout.Offsets[i])
		end := start + int(layout.Lengths[i])
		if i == int(LumpGame) {
			if err := game.CopyTo(raw[start:end]); err != nil {
				return nil, &FormatError{"emit game directory", err}
			}
		} else {
			copy(raw[start:end], file.lumps[i].raw)
		}
	}
	encodeHeader(raw[:HeaderSize], &header, plan.headerLayout)
	return raw, nil
}

// WriteLump returns an owned snapshot of committed stored bytes.
// Game positions remain its stored base; use MarshalGameLump for explicit relocation.
func (w *Writer) WriteLump(index LumpId) ([]byte, error) {
	if w == nil || w.data == nil {
		return nil, fmt.Errorf("%w: nil file", ErrMalformed)
	}
	return w.data.RawLump(index)
}

// NewWriter binds the supplied file without an ownership-by-value facade.
func NewWriter(file *Bsp) *Writer { return &Writer{data: file} }

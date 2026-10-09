package bsp

import (
	"fmt"
	"io"

	"github.com/galaco/bsp/internal/codec"
	"github.com/galaco/bsp/internal/format"
)

type exportPlan struct {
	layout       codec.Layout
	headerLayout format.HeaderLayout
	game         codec.GameRelocation
}

// prepareExport is shared by emission and current-state map checksums.
func (file *File) prepareExport() (exportPlan, error) {
	if file == nil || file.header.ID != fileIdent {
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

// exportHeader derives output descriptors without changing historical metadata.
func (file *File) exportHeader(plan exportPlan) Header {
	header := file.header
	for i := range header.Lumps {
		header.Lumps[i].Offset = plan.layout.Offsets[i]
		header.Lumps[i].Length = plan.layout.Lengths[i]
	}
	return header
}

// MarshalBinary emits a canonical numeric-ID layout with aligned stored payloads.
// Export does not mutate the file; game-directory offsets are relocated.
func (file *File) MarshalBinary() ([]byte, error) {
	plan, err := file.prepareExport()
	if err != nil {
		return nil, err
	}
	layout, game := plan.layout, plan.game
	// Validate the entire layout and game fixups before the output allocation.
	raw := make([]byte, int(layout.Size))
	header := file.exportHeader(plan)
	for i := range header.Lumps {
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

// WriteTo emits the same bytes as MarshalBinary without staging a complete BSP.
// Semantic validation precedes all writes. I/O errors return the actual emitted
// byte count, and the caller retains ownership of the writer.
func (file *File) WriteTo(writer io.Writer) (int64, error) {
	plan, err := file.prepareExport()
	if err != nil {
		return 0, err
	}
	if writer == nil {
		return 0, fmt.Errorf("%w: nil writer", ErrMalformed)
	}
	header := file.exportHeader(plan)
	var headerBytes [HeaderSize]byte
	encodeHeader(headerBytes[:], &header, plan.headerLayout)
	written, err := writeChunk(writer, headerBytes[:])
	if err != nil {
		return written, err
	}
	var zeros [3]byte
	for id, length := range plan.layout.Lengths {
		if length == 0 {
			continue
		}
		padding := int64(plan.layout.Offsets[id]) - written
		if padding != 0 {
			n, err := writeChunk(writer, zeros[:padding])
			written += n
			if err != nil {
				return written, err
			}
		}
		var n int64
		if id == int(LumpGame) {
			n, err = plan.game.WriteTo(writer)
		} else {
			n, err = writeChunk(writer, file.lumps[id].raw)
		}
		written += n
		if err != nil {
			return written, err
		}
	}
	padding := int64(plan.layout.Size) - written
	if padding != 0 {
		n, err := writeChunk(writer, zeros[:padding])
		written += n
		if err != nil {
			return written, err
		}
	}
	return written, nil
}

func writeChunk(writer io.Writer, raw []byte) (int64, error) {
	n, err := writer.Write(raw)
	if n != len(raw) && err == nil {
		err = io.ErrShortWrite
	}
	return int64(n), err
}

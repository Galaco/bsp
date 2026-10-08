package bsp

import "github.com/galaco/bsp/lumps"

// Bsp is the root .bsp filetype container.
// Consists of a 1036byte header and 64 lump blocks.
type Bsp struct {
	header Header
	lumps  [64]Lump
}

// Header is the Bsp header. Contains format and lump layout data.
// Do not trust lump information between import and export
type Header struct {
	Id       int32
	Version  int32
	Lumps    [64]HeaderLump
	Revision int32
}

// HeaderLump contains layout information for a given lump, stored in the Header.
type HeaderLump struct {
	Offset  int32
	Length  int32
	Version int32
	Id      [4]byte
}

// Header gets the header for a bsp.
func (bsp *Bsp) Header() *Header {
	return &bsp.header
}

// Lump decodes a payload and returns its decoding error.
func (bsp *Bsp) Lump(index LumpId) (lumps.ILump, error) {
	lump, err := bsp.RawLump(index)
	if err != nil {
		return nil, err
	}
	return lump.Contents()
}

// RawLump returns a range-checked lump container.
func (bsp *Bsp) RawLump(index LumpId) (*Lump, error) {
	if err := validateID(index); err != nil {
		return nil, err
	}
	return &bsp.lumps[int(index)], nil
}

// SetLump replaces a range-checked lump.
func (bsp *Bsp) SetLump(index LumpId, lump Lump) error {
	if err := validateID(index); err != nil {
		return err
	}
	bsp.lumps[int(index)] = lump
	return nil
}

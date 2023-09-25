package bsp

import "github.com/galaco/bsp/lumps"

// Bsp is the root .bsp filetype container.
// Consists of a 1036byte header and 64 lump blocks.
type Bsp struct {
	header Header
	lumps  [64]lumps.ILump
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

// Lump gets the lump for a given id.
func (bsp *Bsp) Lump(index LumpId) lumps.ILump {
	return bsp.lumps[index]
}

// SetLump sets the lump data for a given id.
func (bsp *Bsp) SetLump(index LumpId, lump lumps.ILump) {
	bsp.lumps[int(index)] = lump
}

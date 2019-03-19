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

// GetHeader gets the header for a bsp.
func (bsp *Bsp) GetHeader() Header {
	return bsp.header
}

// GetLump gets the lump for a given index.
func (bsp *Bsp) GetLump(index LumpId) lumps.ILump {
	return bsp.GetLumpRaw(index).GetContents()
}

// GetLumpRaw gets the lump for a given index.
func (bsp *Bsp) GetLumpRaw(index LumpId) *Lump {
	return &bsp.lumps[int(index)]
}

// SetLump sets the lump data for a given index.
func (bsp *Bsp) SetLump(index LumpId, lump Lump) {
	bsp.lumps[int(index)] = lump
}

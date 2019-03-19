package bsp

import "github.com/galaco/bsp/lumps"

// Root .bsp filetype container.
// Consists of a 1036byte header and 64 lump blocks.
type Bsp struct {
	header Header
	lumps  [64]Lump
}

// Bsp header. Contains format and lump layout data.
// Do not trust lump information between import and export
type Header struct {
	Id       int32
	Version  int32
	Lumps    [64]HeaderLump
	Revision int32
}

// Layout information for a given lump, stored in the Header.
type HeaderLump struct {
	Offset  int32
	Length  int32
	Version int32
	Id      [4]byte
}

// Get the header for a bsp.
func (bsp *Bsp) GetHeader() Header {
	return bsp.header
}

// Get the lump for a given index.
func (bsp *Bsp) GetLump(index LumpId) lumps.ILump {
	return bsp.GetLumpRaw(index).GetContents()
}

// Get the lump for a given index.
func (bsp *Bsp) GetLumpRaw(index LumpId) *Lump {
	return &bsp.lumps[int(index)]
}

// Set the lump data for a given index.
func (bsp *Bsp) SetLump(index LumpId, lump Lump) {
	bsp.lumps[int(index)] = lump
}

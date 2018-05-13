package bsp

// Root .bsp filetype container.
// Consists of a 1036byte header and 64 lump blocks.
type Bsp struct {
	header Header
	lumps [64]Lump
}

// Bsp header. Contains format and lump layout data.
// Do not trust lump information between import and export
type Header struct {
	Id int32
	Version int32
	Lumps [64]HeaderLump
	Revision int32
}

// Layout information for a given lump, stored in the Header.
type HeaderLump struct {
	Offset int32
	Length int32
	Version int32
	Id [4]byte
}

// Get the header for a bsp.
func (bsp *Bsp) GetHeader() Header {
	return bsp.header
}

// Get the lump for a given index.
func (bsp *Bsp) GetLump(index int) *Lump {
	return &bsp.lumps[index]
}

// Set the lump data for a given index.
func (bsp *Bsp) SetLump(index int, lump Lump) {
	bsp.lumps[index] = lump
}
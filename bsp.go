package bsp

// Bsp is the root .bsp filetype container.
// Consists of a 1036byte header and 64 lump blocks.
type Bsp struct {
	Header Header   `json:"header"`
	Lumps  [64]Lump `json:"lumps"`
}

// Header is the Bsp header. Contains format and lump layout data.
// Do not trust lump information between import and export
type Header struct {
	Id       int32          `json:"id"`
	Version  int32          `json:"version"`
	Lumps    [64]HeaderLump `json:"lumps"`
	Revision int32          `json:"revision"`
}

// HeaderLump contains layout information for a given lump, stored in the Header.
type HeaderLump struct {
	// Offset is the offset into the file (in bytes).
	Offset int32 `json:"offset"`
	// Length is the lump length (in bytes).
	Length int32 `json:"length"`
	// Version is the lump version.
	Version int32 `json:"version"`
	// ID is lump ident code.
	// If the lump is compressed then treat as an integer that represents the decompressed size.
	ID [4]byte `json:"id"`
}

// Lump interface.
type Lump interface {
	// FromBytes imports a []byte to a defined lump structure(s).
	FromBytes([]byte) error
	// ToBytes exports lump structure back to []byte.
	ToBytes() ([]byte, error)
	// SetVersion sets bsp version of lump.
	SetVersion(version int32)
}

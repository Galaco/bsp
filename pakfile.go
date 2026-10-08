package bsp

import (
	"io/fs"

	"github.com/galaco/bsp/internal/codec"
	"github.com/galaco/bsp/lumps"
)

// Pakfile is a read-only archive snapshot, independent of later raw replacements.
// It exposes metadata and owned file bytes, not a mutable archive/zip decoder.
// No disk extraction or per-entry editing is provided.
type Pakfile struct{ archive *codec.Pakfile }

// Pakfile returns a read-only view of the current embedded ZIP. An absent lump
// yields an empty view; corrupt nonempty bytes return a contextual parse error.
func (file *File) Pakfile() (*Pakfile, error) {
	decoded, err := file.decodeLump(LumpPakfile)
	if err != nil {
		return nil, err
	}
	return &Pakfile{archive: decoded.(*codec.Pakfile)}, nil
}

// Entries returns metadata in ZIP directory order. Editing it does not edit the
// archive. Use SetRawLump(LumpPakfile, descriptor, rawZIP) for whole replacement.
func (pak *Pakfile) Entries() []lumps.PakEntry {
	if pak == nil {
		return nil
	}
	return pak.archive.Entries()
}

// ReadFile returns owned bytes and always closes the entry reader. Lookup is
// case-insensitive, converts backslashes and collapses repeated separators.
// '..' is literal, not filesystem resolution. First matching directory entry
// wins. Missing names wrap fs.ErrNotExist; empty entries succeed with nil error.
func (pak *Pakfile) ReadFile(name string) ([]byte, error) {
	if pak == nil {
		return nil, &fs.PathError{Op: "read pakfile", Path: name, Err: fs.ErrNotExist}
	}
	return pak.archive.ReadFile(name)
}

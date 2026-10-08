package bsp

import (
	"fmt"
	"github.com/galaco/bsp/internal/versions"
	"github.com/galaco/bsp/lumps"
)

// Lump is a container for a lump. Also includes metadata about the lump.
// N.B. Some information mirrors the header's lump descriptor, but header information should not be trusted after
// import completion.
type Lump struct {
	raw    []byte
	data   lumps.ILump
	length int32
	id     LumpId
	loaded bool
}

// SetId sets lump identifier
// Id is the lump type id (not the id for the order the lumps are stored)
func (l *Lump) SetId(index LumpId) {
	l.id = index
}

// Contents decodes the payload without hiding its underlying error.
func (l *Lump) Contents() (lumps.ILump, error) {
	if l.data == nil {
		return nil, fmt.Errorf("%w: lump %d", ErrUnsupportedLump, l.id)
	}
	if !l.loaded {
		if err := l.data.Unmarshall(l.raw); err != nil {
			return nil, &FormatError{fmt.Sprintf("decode lump %d", l.id), err}
		}
		l.loaded = true
	}
	return l.data, nil
}

// SetContents Set content type of a lump.
func (l *Lump) SetContents(data lumps.ILump) {
	l.data = data
	l.loaded = false
}

// RawContents Get the raw []byte contents of a lump.
// N.B. This is the raw imported value. To get the raw value of a modified lump, use Contents().Marshall()
func (l *Lump) RawContents() []byte {
	return l.raw
}

// SetRawContents Set raw []byte contents of a lump.
func (l *Lump) SetRawContents(raw []byte) {
	l.raw = raw
}

// Length Get length of a lump in bytes.
func (l *Lump) Length() int32 {
	return l.length
}

// getReferenceLumpByIndex Return an instance of a Lump for a given offset.
func getReferenceLumpByIndex(index int, version int32) (lumps.ILump, error) {
	if index < 0 || index > 63 {
		return nil, fmt.Errorf("invalid lump id: %d provided", index)
	}

	l, err := versions.GetLumpForVersion(int(version), index)
	if err != nil {
		return nil, err
	}

	l.SetVersion(version)

	return l, nil
}

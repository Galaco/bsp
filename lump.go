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

// SetId Get lump identifier
// Id is the lump type id (not the id for the order the lumps are stored)
func (l *Lump) SetId(index LumpId) {
	l.id = index
}

// GetContents Get the contents of a lump.
// NOTE: Will need to be cast to the relevant lumps
func (l *Lump) GetContents() lumps.ILump {
	if !l.loaded {
		if l.data.Unmarshall(l.raw) != nil {
			return nil
		}
		l.loaded = true
	}
	return l.data
}

// SetContents Set content type of a lump.
func (l *Lump) SetContents(data lumps.ILump) {
	l.data = data
	l.loaded = false
}

// GetRawContents Get the raw []byte contents of a lump.
// N.B. This is the raw imported value. To get the raw value of a modified lump, use GetContents().Marshall()
func (l *Lump) GetRawContents() []byte {
	return l.raw
}

// SetRawContents Set raw []byte contents of a lump.
func (l *Lump) SetRawContents(raw []byte) {
	l.raw = raw
}

// GetLength Get length of a lump in bytes.
func (l *Lump) GetLength() int32 {
	return l.length
}

// getReferenceLumpByIndex Return an instance of a Lump for a given offset.
func getReferenceLumpByIndex(index int, version int32) (lumps.ILump, error) {
	if index < 0 || index > 63 {
		return nil, fmt.Errorf("invalid lump id: %d provided", index)
	}

	return versions.GetLumpForVersion(int(version), index)
}

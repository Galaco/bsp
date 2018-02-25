package bsp

import (
	"github.com/galaco/bsp/lumps"
	"log"
	"github.com/galaco/bsp/versions"
)

// Container for a lump. Also includes metadata about the lump.
// N.B. Some information mirrors the header's lump descriptor, but header information should not be trusted after
// import completion.
type Lump struct {
	raw []byte
	data lumps.ILump
	length int32
	index int
	loaded bool
}

// Get lump identifier
// Id is the lump type index (not the index for the order the lumps are stored)
func (l *Lump) SetId(index int) {
	l.index = index
}

// Get the contents of a lump.
// NOTE: Will need to be cast to the relevant lumps
func (l *Lump) GetContents() lumps.ILump {
	if l.loaded == false {
		l.data = l.data.FromBytes(l.raw, int32(len(l.raw)))
		l.loaded = true
	}
	return l.data
}

// Set contents of a lump.
func (l *Lump) SetContents(data lumps.ILump) {
	l.data = data
	l.loaded = false
}

// Get the raw []byte contents of a lump.
// N.B. This is the raw imported value. To get the raw value of a modified lump, use GetContents().ToBytes()
func (l *Lump) GetRawContents() []byte {
	return l.raw
}

// Set raw []byte contents of a lump.
func (l *Lump) SetRawContents(raw []byte) {
	l.raw = raw
}

// Get length of a lump in bytes.
func (l *Lump) GetLength() int32 {
	return l.length
}

// Return an instance of a Lump for a given offset.
func getLumpForIndex(index int, version int32) lumps.ILump {
	if index < 0 || index > 63 {
		log.Fatalf("Invalid lump index: %d provided\n", index)
	}

	switch version {
	case 20:
		return versions.GetVersion20Mapping()[index]
	default:
		log.Fatalf("Bsp version: %d not currently supported\n", version)
	}
	return nil
}
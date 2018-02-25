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
		log.Fatal("Invalid lump index provided.")
	}

	switch version {
	case 20:
		return versions.GetVersion20Mapping()[index]
	default:
		log.Fatal("Bsp version not currently supported")
	}
	return nil
}
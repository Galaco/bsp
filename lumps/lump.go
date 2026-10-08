package lumps

import "github.com/galaco/bsp/internal/format"

// Generic is a raw byte lump
// It should be used in cases where the data has no or unknown
// structures.
type Generic struct {
	Metadata
	data []byte
}

// Unmarshall Imports this lump from raw byte data
func (lump *Generic) Unmarshall(data []byte) (err error) {
	lump.length = len(data)
	lump.data = data

	return err
}

// Marshall Dumps this lump back to raw byte data
func (lump *Generic) Marshall() ([]byte, error) {
	return lump.data, nil
}

// Metadata is a Helper info for a lump
type Metadata struct {
	length  int
	context format.Context
}

// Length Returns lump import length in bytes.
func (info *Metadata) Length() int {
	return info.length
}

// SetLength sets lump import length in bytes
func (info *Metadata) SetLength(length int) {
	info.length = length
}

// Context returns the file profile and descriptor version used by this codec.
func (info *Metadata) Context() format.Context {
	return info.context
}

// SetContext sets the file profile and descriptor version used by this codec.
func (info *Metadata) SetContext(context format.Context) {
	info.context = context
}

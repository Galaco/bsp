package lumps

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
	version int32
}

// Length Returns lump import length in bytes.
func (info *Metadata) Length() int {
	return info.length
}

// SetLength sets lump import length in bytes
func (info *Metadata) SetLength(length int) {
	info.length = length
}

// Version Returns lump import version in bytes.
func (info *Metadata) Version() int32 {
	return info.version
}

// SetVersion sets bsp version of lump
func (info *Metadata) SetVersion(version int32) {
	info.version = version
}

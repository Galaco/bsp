package lumps

// LumpGeneric is a raw byte lump
// It should be used in cases where the data has no or unknown
// structures.
type LumpGeneric struct {
	LumpInfo
	data []byte
}

// Unmarshall Imports this lump from raw byte data
func (lump *LumpGeneric) Unmarshall(data []byte, length int32) {
	lump.length = length
	lump.data = data
}

// Marshall Dumps this lump back to raw byte data
func (lump *LumpGeneric) Marshall() ([]byte,error) {
	return lump.data, nil
}

// LumpInfo is a Helper info for a lump
type LumpInfo struct {
	length  int32
	version int32
}

// GetLength Returns lump import length in bytes.
func (info LumpInfo) GetLength() int32 {
	return info.length
}

// SetLength sets lump import length in bytes
func (info LumpInfo) SetLength(length int32) {
	info.length = length
}

// SetVersion sets bsp version of lump
func (info LumpInfo) SetVersion(version int32) {
	info.version = version
}

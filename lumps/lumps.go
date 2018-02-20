package lumps

/**
	Lump interface.
	Organise Lump data in a cleaner and more accessible manner
 */
type ILump interface {
	FromBytes([]byte, int32) ILump

	GetData() interface{}

	ToBytes() []byte
}

/**
	Helper info for a lump
 */
type LumpInfo struct {
	length int32
}
func (info LumpInfo) GetLength() int32 {
	return info.length
}
func (info LumpInfo) SetLength(length int32) {
	info.length = length
}
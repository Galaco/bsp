package lumps

/**
	Lump interface.
	Organise Lump data in a cleaner and more accessible manner
 */
type ILump interface {
	/**
		Import a []byte to a defined lump structure(s).
	 */
	FromBytes([]byte, int32) ILump

	/**
		Return populated lump structure(s).
	 */
	GetData() interface{}

	/**
		Export lump structure back to []byte.
	 */
	ToBytes() []byte
}

/**
	Helper info for a lump
 */
type LumpInfo struct {
	length int32
}

/**
	Return lump import length in bytes.
 */
func (info LumpInfo) GetLength() int32 {
	return info.length
}
/**
	Set lump import length in bytes
 */
func (info LumpInfo) SetLength(length int32) {
	info.length = length
}
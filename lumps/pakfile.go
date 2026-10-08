package lumps

// PakEntry is an owned metadata snapshot in ZIP directory order. Name retains
// its exact archive spelling; resource lookup applies separate normalization.
type PakEntry struct {
	Name             string
	Method           uint16
	Flags            uint16
	CRC32            uint32
	CompressedSize   uint64
	UncompressedSize uint64
	Directory        bool
}

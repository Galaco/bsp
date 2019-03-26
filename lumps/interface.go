package lumps

// ILump Lump interface.
// Organise Lump data in a cleaner and more accessible manner
type ILump interface {
	// Unmarshall Imports a []byte to a defined lump structure(s).
	Unmarshall([]byte) error

	// Marshall Exports lump structure back to []byte.
	Marshall() ([]byte, error)
}

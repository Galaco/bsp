package lumps

// ILump Lump interface.
// Organise Lump data in a cleaner and more accessible manner
type ILump interface {
	// Unmarshall Imports a []byte to a defined lump structure(s).
	Unmarshall([]byte) error

	// Marshall Exports lump structure back to []byte.
	Marshall() ([]byte, error)

	SetVersion(version int32)
}

type lump interface {
	Unmarshall(raw []byte) (err error)
}

func LumpDataToLumpType[T lump](l RawBytes) (*T, error) {
	var t T
	if err := t.Unmarshall(l.data); err != nil {
		return nil, err
	}

	return &t, nil
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

// RawBytes is Lump n: RawBytes lump type
// the contents are just raw bytes, left up to the implementer to handle.
type RawBytes struct {
	Metadata
	data []byte
}

// Unmarshall Imports this lump from raw byte data
func (lump *RawBytes) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.data = raw
	lump.Metadata.SetLength(length)

	return err
}

// GetData gets internal format structure data
func (lump *RawBytes) GetData() []byte {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *RawBytes) Marshall() ([]byte, error) {
	return lump.data, nil
}

type Unimplemented = RawBytes

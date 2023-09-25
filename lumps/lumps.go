package lumps

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"
)

// Lump Lump interface.
// Organise Lump data in a cleaner and more accessible manner
type Lump interface {
	// FromBytes imports a []byte to a defined lump structure(s).
	FromBytes([]byte) error
	// ToBytes exports lump structure back to []byte.
	ToBytes() ([]byte, error)
	// SetVersion sets bsp version of lump.
	SetVersion(version int32)
}

// unmarshallBasicLump is a helper function for lumps that are just a single []T.
func unmarshallBasicLump[V any](raw []byte) (Metadata, []V, error) {
	var meta Metadata
	length := len(raw)
	meta.SetLength(length)
	if length == 0 {
		return meta, nil, fmt.Errorf("lump is empty")
	}

	var sampleV V
	v := make([]V, length/int(unsafe.Sizeof(sampleV)))
	if err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &v); err != nil {
		return meta, nil, err
	}
	meta.SetLength(length)

	return meta, v, nil
}

func marshallBasicLump(data any) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, data)
	return buf.Bytes(), err
}

// Metadata is a Helper info for a lump
type Metadata struct {
	length  int
	version int32
}

func (info *Metadata) GetMetadata() *Metadata {
	return info
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

// FromBytes imports this lump from raw byte data
func (lump *RawBytes) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.data = raw
	lump.Metadata.SetLength(length)

	return err
}

// Contents returns internal format structure data
func (lump *RawBytes) Contents() []byte {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *RawBytes) ToBytes() ([]byte, error) {
	return lump.data, nil
}

type Unimplemented = RawBytes

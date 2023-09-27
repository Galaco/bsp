package lump

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

// unmarshallBasicLump is a helper function for unmarshalling []byte to lumps that are just a single []T.
func unmarshallBasicLump[V any](raw []byte) (Metadata, []V, error) {
	var meta Metadata
	if len(raw) == 0 {
		return meta, nil, nil
	}

	var sampleV V
	v := make([]V, len(raw)/int(unsafe.Sizeof(sampleV)))
	if err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &v); err != nil {
		return meta, nil, err
	}

	return meta, v, nil
}

// marshallBasicLump is a helper function for marshalling lumps that are just a single []T to []byte.
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
	Data []byte `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *RawBytes) FromBytes(raw []byte) error {
	lump.Data = raw

	return nil
}

// Contents returns internal format structure Data
func (lump *RawBytes) Contents() []byte {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *RawBytes) ToBytes() ([]byte, error) {
	return lump.Data, nil
}

type Unimplemented = RawBytes

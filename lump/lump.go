package lump

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

// Unimplemented is Lump n: Unimplemented lump type.
type Unimplemented = rawBytes

// Metadata is a Helper info for a lump
type Metadata struct {
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

// rawBytes is Lump n: rawBytes lump type
// the contents are just raw bytes, left up to the implementer to handle.
type rawBytes struct {
	Metadata
	Data []byte `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *rawBytes) FromBytes(raw []byte) error {
	lump.Data = raw

	return nil
}

// Contents returns internal format structure Data
func (lump *rawBytes) Contents() []byte {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *rawBytes) ToBytes() ([]byte, error) {
	return lump.Data, nil
}

// unmarshallBasicLump is a helper function for unmarshalling []byte to lumps that are just a single []T.
func unmarshallBasicLump[V any](raw []byte) ([]V, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	var sampleV V
	v := make([]V, len(raw)/int(unsafe.Sizeof(sampleV)))
	if err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// marshallBasicLump is a helper function for marshalling lumps that are just a single []T to []byte.
func marshallBasicLump(data any) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, data)
	return buf.Bytes(), err
}

// unmarshallTaggedLump is a helper function for unmarshalling []byte to lumps that are just a single []T.
// It uses the "bsp" tag on struct fields to determine which fields to read.
// It is generally useful for lumps that have multiple versions and have different fields & lengths of T.
func unmarshallTaggedLump[T any](raw []byte, version string) ([]T, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	var sampleT T

	// Figure out the length of our struct for our version.
	var binarylenT int
	for _, field := range reflect.VisibleFields(reflect.TypeOf(sampleT)) {
		// Fields that are in this version should contribute to the length.
		if t := field.Tag.Get("bsp"); t == "" || t == version {
			binarylenT += int(field.Type.Size())
		}
	}
	if len(raw)%binarylenT != 0 {
		// length doesn't match exactly a multiple of our calculated struct size.
		return nil, fmt.Errorf("lump length %d is not a multiple of %d", len(raw), binarylenT)
	}

	// Calculate the padding size to properly read the struct.
	padSize := int(unsafe.Sizeof(sampleT)) - binarylenT

	v := make([]T, len(raw)/binarylenT)
	for i := range v {
		b2 := append([]byte{}, raw[(binarylenT*i):(binarylenT*i)+binarylenT]...)
		buf := bytes.NewBuffer(append(b2, []byte(strings.Repeat("\x00", padSize))...))
		if err := binary.Read(buf, binary.LittleEndian, &v[i]); err != nil {
			return nil, err
		}
	}

	return v, nil
}

// marshallTaggedLump is a helper function for marshalling lumps that are just a single []T to []byte.
// It uses the "bsp" tag on struct fields to determine which fields to write.
// It is generally useful for lumps that have multiple versions and have different fields & lengths of V.
func marshallTaggedLump[T any](data []T, version string) ([]byte, error) {
	if len(data) == 0 {
		return []byte{}, nil
	}

	var buf bytes.Buffer
	visibleFields := reflect.VisibleFields(reflect.TypeOf(data[0]))
	for _, v := range data {
		for _, field := range visibleFields {
			rv := reflect.Indirect(reflect.ValueOf(&v))
			// Fields that are in this version should contribute to the length.
			if t := field.Tag.Get("bsp"); t == "" || t == version {
				switch rv.FieldByName(field.Name).Kind() {
				//case reflect.Slice:
				//	if err := binary.Write(&buf, binary.LittleEndian, rv.FieldByName(field.Name).Elem()); err != nil {
				//		return nil, err
				//	}
				default:
					if err := binary.Write(&buf, binary.LittleEndian, rv.FieldByName(field.Name).Interface()); err != nil {
						return nil, err
					}
				}
			}
		}
	}
	return buf.Bytes(), nil
}

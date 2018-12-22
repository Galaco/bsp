package lumps

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"strings"
)

// Lump 40: Pakfile
type Pakfile struct {
	LumpGeneric
	zipReader *zip.Reader
}

// Import this lump from raw byte data
func (lump *Pakfile) FromBytes(raw []byte, length int32) {
	lump.data = raw
	lump.LumpInfo.SetLength(length)

	b := bytes.NewReader(raw)
	zipReader, err := zip.NewReader(b, int64(length))
	if err == nil {
		lump.zipReader = zipReader
	}
}

// GetData Get internal format structure data
func (lump *Pakfile) GetData() *zip.Reader {
	return lump.zipReader
}

// ToBytes Returns the contents of this lump as a []byte
func (lump *Pakfile) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

// GetFile Get a specific file from the pak
// This function does have the expectation that filenames are case-insensitive,
// so in the (shouldn't happen) case of name collisions over case-insensitivity,
// there may be issues
func (lump *Pakfile) GetFile(filePath string) ([]byte, error) {
	filePath = lump.sanitisePath(filePath)
	for _, f := range lump.zipReader.File {
		if lump.sanitisePath(f.Name) == filePath {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			return ioutil.ReadAll(rc)
		}
	}
	return []byte{}, nil
}

// sanitisePath ensures that the requested path matches internal casing
func (lump *Pakfile) sanitisePath(filePath string) string {
	return strings.Replace(strings.ToLower(filePath), "\\", "/", -1)
}

package lumps

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"strings"
)

// Pakfile is Lump 40: Pakfile
type Pakfile struct {
	Generic
	zipReader *zip.Reader
}

// Unmarshall Imports this lump from raw byte data
func (lump *Pakfile) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.data = raw
	lump.Metadata.SetLength(length)

	b := bytes.NewReader(raw)
	zipReader, err := zip.NewReader(b, int64(length))
	if err == nil {
		lump.zipReader = zipReader
	}

	return err
}

// GetData GetData gets internal format structure data
func (lump *Pakfile) GetData() *zip.Reader {
	return lump.zipReader
}

// Marshall Returns the contents of this lump as a []byte
func (lump *Pakfile) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}

// GetFile Get a specific file from the pak
// This function does have the expectation that filenames are case-insensitive,
// so in the (shouldn't happen) case of name collisions over case-insensitivity,
// there may be issues
func (lump *Pakfile) GetFile(filePath string) ([]byte, error) {
	filePath = lump.sanitisePath(filePath)
	for _, f := range lump.zipReader.File {
		if strings.EqualFold(filePath, lump.sanitisePath(f.Name)) {
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
	return strings.Replace(
		strings.Replace(strings.ToLower(filePath), "\\", "/", -1),
		"//",
		"/",
		-1)
}

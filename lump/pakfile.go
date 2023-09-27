package lump

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
)

// Pakfile is Lump 40: Pakfile
type Pakfile struct {
	RawBytes
	zipReader *zip.Reader
}

// FromBytes imports this lump from raw byte Data
func (lump *Pakfile) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Data = raw
	lump.Metadata.SetLength(length)

	b := bytes.NewReader(lump.Data)
	zipReader, err := zip.NewReader(b, int64(length))
	if err != nil {
		return err
	}
	lump.zipReader = zipReader
	return nil
}

// Contents gets internal format structure Data
func (lump *Pakfile) Contents() *zip.Reader {
	return lump.zipReader
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
			return io.ReadAll(rc)
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

package lumps

import (
	"encoding/binary"
	"bytes"
	"archive/zip"
	"io/ioutil"
	"strings"
)

/**
	Lump 40: Pakfile
 */
type Pakfile struct {
	LumpGeneric
	zipReader *zip.Reader
}

func (lump *Pakfile) FromBytes(raw []byte, length int32) {
	lump.data = raw
	lump.LumpInfo.SetLength(length)

	b := bytes.NewReader(raw)
	zipReader,err := zip.NewReader(b, int64(length))
	if err == nil {
		lump.zipReader = zipReader
	}
}

// Returns this lumps contents as whatever internal format this
// lump uses (requires casting)
func (lump *Pakfile) GetData() *zip.Reader {
	return lump.zipReader
}

// Returns the contents of this lump as a []byte
func (lump *Pakfile) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}

// Get a specific file from the pak
func (lump *Pakfile) GetFile(filePath string) ([]byte,error) {
	filePath = strings.ToLower(filePath)
	for _,f := range lump.zipReader.File {
		if strings.ToLower(f.Name) == filePath {
			rc,err := f.Open()
			if err != nil {
				return nil,err
			}
			return ioutil.ReadAll(rc)
		}
	}
	return []byte{},nil
}
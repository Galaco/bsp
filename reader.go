package bsp

import (
	"bytes"
	"github.com/galaco/bsp/lumps"
	"unsafe"
	"io"
	"encoding/binary"
	"os"
)

// Bsp File reader.
type Reader struct {
	stream io.Reader
}

// Reads the BSP into internal byte structure
// Note that parsing is somewhat lazy. Proper data structures are only generated for
// lumps that are requested at a later time. This generated the header, then []byte
// data for each lump
func (r *Reader) Read() (*Bsp,error) {
	bsp := Bsp{}

	buf := bytes.Buffer{}
	_,err := buf.ReadFrom(r.stream)
	if err != nil {
		return nil,err
	}
	reader := bytes.NewReader(buf.Bytes())

	//Create Header
	h,err := r.readHeader(reader, bsp.header)
	if err != nil {
		return nil,err
	}
	bsp.header = *h

	// Create lumps from header data
	for index := range bsp.header.Lumps {
		lp,err := r.readLump(reader, bsp.header, index)
		if err != nil {
			return nil,err
		}
		bsp.lumps[index].SetId(index)
		bsp.lumps[index].SetRawContents(lp)
		refLump,err := getReferenceLumpByIndex(index, bsp.header.Version)

		// There are specific rules for the game lump that requires some extra information
		// Game lump lumps have offset data relative to file start, not lump start
		// This will correct the offsets to the start of the lump.
		// @NOTE: Portal2 console uses relative offsets. This game+platform are not supported currently
		if index == LUMP_GAME_LUMP {
			refLump.(*lumps.Game).UpdateInternalOffsets(bsp.header.Lumps[index].Offset)
		}

		if err != nil {
			return nil,err
		}
		bsp.lumps[index].SetContents(refLump)
	}

	return &bsp,err
}

// Parse header from the bsp file.
func (r *Reader) readHeader(reader *bytes.Reader, header Header) (*Header, error) {
	headerSize := unsafe.Sizeof(header)
	headerBytes := make([]byte, headerSize)

	sectionReader := io.NewSectionReader(reader, 0, int64(len(headerBytes)))
	_, err := sectionReader.Read(headerBytes)
	if err != nil {
		return nil, err
	}

	err = binary.Read(bytes.NewBuffer(headerBytes[:]), binary.LittleEndian, &header)
	if err != nil {
		return nil,err
	}

	return &header,nil
}

// Reads a single lumps data
// Expect a byte reader containing the lump data, as well as the
// header and lump identifier (index)
func (r *Reader) readLump(reader *bytes.Reader, header Header, index int) ([]byte, error) {
	//Limit lump data to declared size
	lumpHeader := header.Lumps[index]
	raw := make([]byte, lumpHeader.Length)

	// Skip reading for empty lump
	if lumpHeader.Length > 0 {
		sectionReader := io.NewSectionReader(reader, int64(lumpHeader.Offset), int64(lumpHeader.Length))
		_, err := sectionReader.Read(raw)
		if err != nil {
			return nil, err
		}
	}

	return raw,nil
}

// Wraps ReadFromStream to control the file access as well.
// Use ReadFromStream if you already have a file handle
func ReadFromFile(filepath string) (*Bsp, error) {
	f, err := os.Open(filepath)
	defer f.Close()
	if err != nil {
		return nil,err
	}
	b,err := ReadFromStream(f)

	return b,err
}

// Read from any struct that implements io.Reader
// handy for passing in a string/bytes/other stream
func ReadFromStream(reader io.Reader) (*Bsp, error) {
	r := &Reader {
		reader,
	}

	return r.Read()
}
package bsp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/galaco/bsp/lumps"
	"io"
	"os"
	"unsafe"
)

// Reader is a Bsp File reader.
type Reader struct {
	stream io.Reader
}

// Read reads the BSP into internal byte structure
// Note that parsing is somewhat lazy. Proper data structures are only generated for
// lumps that are requested at a later time. This generated the header, then []byte
// data for each lump
func (r *Reader) Read() (bsp *Bsp, err error) {
	defer func() {
		if r := recover(); r != nil {
			bsp = nil
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
		}
	}()

	bsp = &Bsp{}

	buf := bytes.Buffer{}
	_, err = buf.ReadFrom(r.stream)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(buf.Bytes())

	//Create Header
	h, err := r.readHeader(reader, bsp.header)
	if err != nil {
		return nil, err
	}
	bsp.header = *h

	// Create lumps from header data
	for index := range bsp.header.Lumps {
		lp, err := r.readLump(reader, bsp.header, index)
		if err != nil {
			return nil, err
		}
		bsp.lumps[index].SetId(LumpId(index))
		bsp.lumps[index].SetRawContents(lp)
		refLump, err := getReferenceLumpByIndex(index, bsp.header.Version)

		// There are specific rules for the game lump that requires some extra information
		// Game lump lumps have offset data relative to file start, not lump start
		// This will correct the offsets to the start of the lump.
		// @NOTE: Portal2 console uses relative offsets. This game+platform are not supported currently
		if index == int(LumpGame) {
			refLump.(*lumps.Game).UpdateInternalOffsets(bsp.header.Lumps[index].Offset)
		}

		if err != nil {
			return nil, err
		}
		bsp.lumps[index].SetContents(refLump)
	}

	return bsp, err
}

// readHeader Parses header from the bsp file.
func (r *Reader) readHeader(reader *bytes.Reader, header Header) (*Header, error) {
	headerSize := unsafe.Sizeof(header)
	headerBytes := make([]byte, headerSize)

	sectionReader := io.NewSectionReader(reader, 0, int64(len(headerBytes)))
	_, err := sectionReader.Read(headerBytes)
	if err != nil {
		return nil, err
	}

	err = binary.Read(bytes.NewBuffer(headerBytes), binary.LittleEndian, &header)
	if err != nil {
		return nil, err
	}

	return &header, nil
}

// readLump Reads a single lumps data
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

	return raw, nil
}

// ReadFromFile Wraps ReadFromStream to control the file access as well.
// Use ReadFromStream if you already have a file handle
func ReadFromFile(filepath string) (*Bsp, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	b, err := ReadFromStream(f)
	if err != nil {
		err2 := f.Close()
		if err2 != nil {
			return nil, err2
		}
		return nil, err
	}

	err = f.Close()
	return b, err
}

// ReadFromStream Reads from any struct that implements io.Reader
// handy for passing in a string/bytes/other stream
func ReadFromStream(reader io.Reader) (*Bsp, error) {
	r := &Reader{
		reader,
	}

	return r.Read()
}

package bsp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"unsafe"

	"github.com/galaco/bsp/lumps"
)

// Reader is a Bsp File reader.
type Reader struct {
}

// Read reads the BSP into internal byte structure
// Note that parsing is somewhat lazy. Proper data structures are only generated for
// lumps that are requested at a later time. This generated the header, then []byte
// data for each lump
func (r *Reader) Read(stream io.Reader) (bsp *Bsp, err error) {
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
	_, err = buf.ReadFrom(stream)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(buf.Bytes())

	//Create Header
	h, err := r.readHeader(reader)
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
		refLump, err := getReferenceLumpByIndex(index, bsp.header.Version)
		if err != nil {
			return nil, err
		}

		// There are specific rules for the game lump that requires some extra information
		// Game lump lumps have offset data relative to file start, not lump start
		// This will correct the offsets to the start of the lump.
		// @NOTE: Portal2 console uses relative offsets. This game+platform are not supported currently
		if index == int(LumpGame) {
			refLump.(*lumps.Game).UpdateInternalOffsets(bsp.header.Lumps[index].Offset)
		}

		if err := refLump.FromBytes(lp); err != nil {
			return nil, err
		}
		bsp.lumps[index] = refLump
	}

	return bsp, err
}

// readHeader Parses header from the bsp file.
func (r *Reader) readHeader(reader *bytes.Reader) (*Header, error) {
	var header Header
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
// header and lump identifier (id)
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

// ReadFromStream Reads from any struct that implements io.Reader
// handy for passing in a string/bytes/other stream
func ReadFromStream(stream io.Reader) (*Bsp, error) {
	r := &Reader{}

	return r.Read(stream)
}

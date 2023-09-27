package bsp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"unsafe"

	"github.com/galaco/bsp/lump"
)

// Reader is a Bsp File reader.
type Reader struct {
	lumpResolver LumpResolver
}

// NewReader creates a new Bsp reader.
func NewReader(lumpResolver LumpResolver) *Reader {
	return &Reader{
		lumpResolver: lumpResolver,
	}
}

// Read reads from an io.Reader into a structured Bsp.
func (r *Reader) Read(stream io.Reader) (bsp *Bsp, err error) {
	// Safeguard against panics.
	// @TODO this needs to be removed and the library made safe.
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

	buf := bytes.NewBuffer([]byte{})
	if _, err := io.Copy(buf, stream); err != nil {
		return nil, err
	}
	reader := bytes.NewReader(buf.Bytes())

	bsp = &Bsp{}

	// Create Header
	header, err := r.readHeader(reader)
	if err != nil {
		return nil, err
	}
	bsp.Header = header

	// Create lumps from header data
	for index := range bsp.Header.Lumps {
		lp, err := r.readLump(reader, &bsp.Header.Lumps[index])
		if err != nil {
			return nil, err
		}
		refLump, err := r.lumpResolver(index, int(bsp.Header.Version))
		if err != nil {
			return nil, err
		}

		// There are specific rules for the game lump that requires some extra information
		// Game lump lumps have offset data relative to file start, not lump start
		// This will correct the offsets to the start of the lump.
		// @NOTE: Portal2 console uses relative offsets. This game+platform are not supported currently
		if index == int(LumpGame) {
			refLump.(*lump.Game).SetAbsoluteFileOffset(int(header.Lumps[index].Offset))
		}

		if err := refLump.FromBytes(lp); err != nil {
			return nil, err
		}
		bsp.Lumps[index] = refLump
	}

	return bsp, err
}

// readHeader Parses header from the bsp file.
func (r *Reader) readHeader(reader io.ReaderAt) (header Header, err error) {
	headerBytes := make([]byte, unsafe.Sizeof(header))

	sectionReader := io.NewSectionReader(reader, 0, int64(len(headerBytes)))
	if _, err := sectionReader.Read(headerBytes); err != nil {
		return header, err
	}

	if err := binary.Read(bytes.NewBuffer(headerBytes), binary.LittleEndian, &header); err != nil {
		return header, err
	}

	return header, nil
}

// readLump reads a single lumps data,
// Expects a byte reader containing the lump data, as well as the
// header and lump identifier (id).
func (r *Reader) readLump(reader *bytes.Reader, headerLump *HeaderLump) ([]byte, error) {
	// Limit lump data to declared size.
	raw := make([]byte, headerLump.Length)

	// Skip reading for empty lump.
	if headerLump.Length == 0 {
		return raw, nil
	}

	if _, err := io.NewSectionReader(reader, int64(headerLump.Offset), int64(headerLump.Length)).Read(raw); err != nil {
		return nil, err
	}

	return raw, nil
}

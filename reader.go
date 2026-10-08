package bsp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/galaco/bsp/lumps"
)

// Reader is a Bsp File reader.
type Reader struct {
	stream io.Reader
}

// Read buffers a validated BSP container; typed payloads are decoded on access.
func (r *Reader) Read() (*Bsp, error) {
	if r.stream == nil {
		return nil, fmt.Errorf("%w: nil reader", ErrMalformed)
	}
	raw, err := io.ReadAll(r.stream)
	if err != nil {
		return nil, err
	}
	if len(raw) < HeaderSize {
		return nil, &FormatError{"header", io.ErrUnexpectedEOF}
	}
	file := &Bsp{}
	if err := binary.Read(bytes.NewReader(raw[:HeaderSize]), binary.LittleEndian, &file.header); err != nil {
		return nil, &FormatError{"header", err}
	}
	if file.header.Id == 0x56425350 {
		return nil, fmt.Errorf("%w: big-endian BSP", ErrUnsupportedFormat)
	}
	if file.header.Id != fileIdent {
		return nil, fmt.Errorf("%w: invalid magic", ErrMalformed)
	}
	if file.header.Version < 19 || file.header.Version > 21 {
		return nil, fmt.Errorf("%w: version %d", ErrUnsupportedFormat, file.header.Version)
	}
	for index, descriptor := range file.header.Lumps {
		start, end, err := checkedLumpRange(descriptor.Offset, descriptor.Length, len(raw))
		if err != nil {
			return nil, &FormatError{fmt.Sprintf("lump %d offset %d length %d", index, descriptor.Offset, descriptor.Length), err}
		}
		decoder, err := getReferenceLumpByIndex(index, file.header.Version)
		if err != nil {
			return nil, err
		}
		if game, ok := decoder.(*lumps.Game); ok {
			game.UpdateInternalOffsets(descriptor.Offset)
		}
		file.lumps[index].SetId(LumpId(index))
		file.lumps[index].SetRawContents(raw[start:end])
		file.lumps[index].SetContents(decoder)
	}
	return file, nil
}

func checkedLumpRange(offset, length int32, total int) (int, int, error) {
	if offset < 0 || length < 0 {
		return 0, 0, fmt.Errorf("%w: negative range", ErrMalformed)
	}
	if length == 0 {
		return 0, 0, nil
	}
	end := int64(offset) + int64(length)
	if offset < HeaderSize {
		return 0, 0, fmt.Errorf("%w: payload overlaps header", ErrMalformed)
	}
	if end > int64(total) {
		return 0, 0, io.ErrUnexpectedEOF
	}
	return int(offset), int(end), nil
}

// ReadFromFile Wraps ReadFromStream to control the file access as well.
// Use ReadFromStream if you already have a file handle
func ReadFromFile(filepath string) (*Bsp, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	b, readErr := ReadFromStream(f)
	closeErr := f.Close()
	if readErr != nil {
		return nil, readErr
	}
	return b, closeErr
}

// ReadFromStream Reads from any struct that implements io.Reader
// handy for passing in a string/bytes/other stream
func ReadFromStream(reader io.Reader) (*Bsp, error) {
	r := &Reader{
		reader,
	}

	return r.Read()
}

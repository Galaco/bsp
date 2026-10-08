package bsp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// Read validates the dialect before buffering payloads; the caller owns the stream.
func Read(reader io.Reader, options ReadOptions) (*File, error) {
	if reader == nil {
		return nil, fmt.Errorf("%w: nil reader", ErrMalformed)
	}
	limits, err := options.limits()
	if err != nil {
		return nil, err
	}
	var headerBytes [HeaderSize]byte
	if _, err := io.ReadFull(reader, headerBytes[:]); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return nil, &FormatError{"header", err}
	}
	magic := binary.LittleEndian.Uint32(headerBytes[:])
	if magic == 0x56425350 {
		return nil, fmt.Errorf("%w: big-endian BSP", ErrUnsupportedFormat)
	}
	if magic != uint32(fileIdent) {
		return nil, fmt.Errorf("%w: invalid magic", ErrMalformed)
	}
	version := int32(binary.LittleEndian.Uint32(headerBytes[4:]))
	layout, err := profileLayout(options.Profile, version)
	if err != nil {
		return nil, err
	}
	file := &File{header: decodeHeader(headerBytes[:], layout), profile: options.Profile, layout: layout, limits: limits}
	var input bytes.Buffer
	input.Grow(HeaderSize)
	input.Write(headerBytes[:])
	if _, err := input.ReadFrom(reader); err != nil {
		return nil, err
	}
	raw := input.Bytes()
	for index, descriptor := range file.header.Lumps {
		start, end, err := checkedLumpRange(descriptor.Offset, descriptor.Length, len(raw))
		if err != nil {
			return nil, &FormatError{fmt.Sprintf("lump %d offset %d length %d", index, descriptor.Offset, descriptor.Length), err}
		}
		file.lumps[index] = payload{raw: raw[start:end:end]}
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

// ReadFile reads a BSP and closes the file it opens.
func ReadFile(filepath string, options ReadOptions) (*File, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	b, readErr := Read(f, options)
	closeErr := f.Close()
	if readErr != nil {
		return nil, readErr
	}
	return b, closeErr
}

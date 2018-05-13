package bsp

import (
	"bytes"
	"unsafe"
	"io"
	"encoding/binary"
)

// Bsp File reader.
type Reader struct {
	stream io.Reader
}

// Parse the set buffer.
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

// Parse a single lump.
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

// Return a new instance of Reader.
func NewReader(reader io.Reader) Reader {
	return Reader{
		reader,
	}
}
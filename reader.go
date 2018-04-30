package bsp

import (
	"bytes"
	"unsafe"
	"io"
	"log"
	"encoding/binary"
)

// Bsp File reader.
type Reader struct {
	stream io.Reader
}

// Parse the set buffer.
func (r *Reader) Read() Bsp {
	bsp := Bsp{}

	buf := bytes.Buffer{}
	_,err := buf.ReadFrom(r.stream)
	if err != nil {
		log.Println(err)
	}
	reader := bytes.NewReader(buf.Bytes())

	//Create Header
	bsp.header = r.readHeader(reader, bsp.header)

	// Create lumps from header data
	for index := range bsp.header.Lumps {
		func(index int, l *Lump) {
			l.SetRawContents(r.readLump(reader, bsp.header, index))
			l.SetContents(getReferenceLumpByIndex(index, bsp.header.Version))
			l.SetId(index)
		}(index, &bsp.lumps[index])
	}

	return bsp
}

// Parse header from the bsp file.
func (r *Reader) readHeader(reader *bytes.Reader, header Header) Header {
	headerSize := unsafe.Sizeof(header)
	headerBytes := make([]byte, headerSize)

	sectionReader := io.NewSectionReader(reader, 0, int64(len(headerBytes)))
	_, err := sectionReader.Read(headerBytes)
	if err != nil {
		log.Fatal(err)
	}

	err = binary.Read(bytes.NewBuffer(headerBytes[:]), binary.LittleEndian, &header)
	if err != nil {
		log.Fatal(err)
	}

	return header
}

// Parse a single lump.
func (r *Reader) readLump(reader *bytes.Reader, header Header, index int) []byte {
	//Limit lump data to declared size
	lumpHeader := header.Lumps[index]
	raw := make([]byte, lumpHeader.Length)

	// Skip reading for empty lump
	if lumpHeader.Length > 0 {
		sectionReader := io.NewSectionReader(reader, int64(lumpHeader.Offset), int64(lumpHeader.Length))
		_, err := sectionReader.Read(raw)
		if err != nil {
			log.Fatal(err)
		}
	}

	return raw
}

// Return a new instance of Reader.
func NewReader(reader io.Reader) Reader {
	return Reader{
		reader,
	}
}
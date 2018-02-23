package bsp

import (
	"bytes"
	"unsafe"
	"io"
	"log"
	"encoding/binary"
	"github.com/galaco/bsp/lumps"
	"fmt"
)

// Bsp File reader.
type Reader struct {
	data []byte
}

// Get buffer to read from.
func (r *Reader) GetBuffer() []byte {
	return r.data
}

// Set buffer to read from.
func (r *Reader) SetBuffer(data []byte) {
	fmt.Println()
	r.data = data
}

// Parse the set buffer.
func (r *Reader) Read() Bsp {
	bsp := Bsp{}

	reader := bytes.NewReader(r.GetBuffer())

	//Create Header
	bsp.header = r.readHeader(reader, bsp.header)

	//Create lumps from header data
	for index := range bsp.header.Lumps {
		bsp.lumps[index].SetContents(r.readLump(reader, bsp.header, index))
	}

	return bsp
}


// Parse header from the bsp file.
func (r Reader) readHeader(reader *bytes.Reader, header Header) Header {
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
func (r Reader) readLump(reader *bytes.Reader, header Header, index int) lumps.ILump{
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

	lump := getLumpForIndex(index, header.Version).FromBytes(raw, lumpHeader.Length)

	return lump
	/*result := lumpData[index].ToBytes()
	fmt.Println(index, len(raw), len(result))
	for i := range raw {
		if raw[i] != result[i] {
			fmt.Println(i, raw[i], result[i])
		}
	}*/
}

// Return a new instance of Reader.
func NewReader() Reader {
	return Reader{}
}
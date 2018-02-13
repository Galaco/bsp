package bsp

import (
	"bytes"
	"io"
	"os"
	"unsafe"
	"log"
	"encoding/binary"
	"github.com/galaco/bsp/lumps"
)

type Bsp struct {
	Header Header
	Lumps [64]lumps.ILump
}

type Header struct {
	Id int32
	Version int32
	Lumps [64]HeaderLump
	Revision int32
}

type HeaderLump struct {
	Offset int32
	Length int32
	Version int32
	Id [4]byte
}

/**
	Parse a bsp and return a structure.
 */
func Parse(file *os.File) Bsp {
	bsp := Bsp{}

	fi,err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	//Load file
	fileData := make([]byte, fi.Size())
	file.Read(fileData)
	reader := bytes.NewReader(fileData)

	//Create Header
	header := Header{}
	header = readHeader(reader, header)
	bsp.Header = header

	//Create lumps from header data
	lumpData := [64]lumps.ILump{}
	lumpData = readLumps(reader, header, lumpData)
	bsp.Lumps = lumpData

	return bsp
}

/**
	Read header from the bsp file.
 */
func readHeader(reader *bytes.Reader, header Header) Header {
	headerSize := unsafe.Sizeof(header)
	headerBytes := make([]byte, headerSize)

	sectionReader := io.NewSectionReader(reader, 0, int64(len(headerBytes)))
	_, err := sectionReader.Read(headerBytes)
	if err != nil {
		log.Fatal(err)
	}

	err = binary.Read(bytes.NewBuffer(headerBytes[:]), binary.LittleEndian, &header)

	return header
}

/**
	Read lumps from bsp file.
 */
func readLumps(reader *bytes.Reader, header Header, lumpData [64]lumps.ILump) [64]lumps.ILump {
	for index,lumpHeader := range header.Lumps {

		//Skip if 0 length b/c nothing to read
		if lumpHeader.Length == 0 {
			continue
		}
		//Limit lump data to declared size
		raw := make([]byte, lumpHeader.Length)
		sectionReader := io.NewSectionReader(reader, int64(lumpHeader.Offset), int64(lumpHeader.Length))


		_, err := sectionReader.Read(raw)
		if err != nil {
			log.Fatal(err)
		}

		lumpData[index] = lumps.GetLumpForIndex(index).FromRaw(raw, lumpHeader.Length)

	}

	return lumpData
}
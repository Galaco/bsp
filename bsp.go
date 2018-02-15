package bsp

import (
	"bytes"
	"io"
	"os"
	"unsafe"
	"log"
	"encoding/binary"
	"github.com/galaco/bsp/lumps"
	"fmt"
	//"reflect"
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

	file.Seek(0,0)

	fi,err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	//Load file
	fileData := make([]byte, fi.Size())
	file.Read(fileData)
	reader := bytes.NewReader(fileData)

	//Create Header
	bsp.Header = readHeader(reader, bsp.Header)

	//Create lumps from header data
	bsp.Lumps = readLumps(reader, bsp.Header, bsp.Lumps)

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
		//Limit lump data to declared size
		raw := make([]byte, lumpHeader.Length)


		// Skip reading for empty lump
		if lumpHeader.Length > 0 {
			sectionReader := io.NewSectionReader(reader, int64(lumpHeader.Offset), int64(lumpHeader.Length))
			_, err := sectionReader.Read(raw)
			if err != nil {
				log.Fatal(err)
			}
		}

		lumpData[index] = lumps.GetLumpForIndex(index).FromBytes(raw, lumpHeader.Length)
	}

	return lumpData
}

/**
	Write Bsp struct to []byte for export.
 */
func ToBytes(bsp Bsp) []byte {
	// First we need to update the header to reflect any lump changes
	// At the same time we can dump our lumps as bytes to write later
	lumpBytes := make([][]byte, 64)
	currentOffset := 1036 // Header always 1036bytes

	for index,lump := range bsp.Lumps {
		// We have to handle lump 35 (GameData differently)
		// Because valve mis-designed the file format and relatively positioned data contains absolute file offsets.
		// NOTE: There is no valid reason for this, as the lump can be located anywhere.
		if index == 35 {
			// @TODO Not yet handled
			gamelump := (bsp.Lumps[index]).(lumps.Game)
			gamelump.UpdateInternalOffsets(int32(currentOffset) - bsp.Header.Lumps[index].Offset)
			lumpBytes[index] = gamelump.ToBytes()
		} else {
			lumpBytes[index] = lump.ToBytes()
		}

		lumpSize := len(lumpBytes[index])

		fmt.Println(index, bsp.Header.Lumps[index].Length, lumpSize)


		// Note these may not change, but we trust our new data more than the header
		bsp.Header.Lumps[index].Length = int32(lumpSize)
		bsp.Header.Lumps[index].Offset = int32(currentOffset)

		currentOffset += lumpSize
	}


	// Now we can export our bsp
	var buf bytes.Buffer

	//Write Header
	binary.Write(&buf, binary.LittleEndian, bsp.Header)

	//Write lumps
	for _,lumpData := range lumpBytes {
		binary.Write(&buf, binary.LittleEndian, lumpData)

	}

	return buf.Bytes()
}
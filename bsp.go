package bsp

import (
	"bytes"
	"io"
	"os"
	"unsafe"
	"log"
	"encoding/binary"
	"github.com/galaco/bsp/lumps"
	"github.com/galaco/bsp/versions"
)

type Bsp struct {
	header Header
	lumps [64]lumps.ILump
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

func (bsp Bsp) GetLump(index int) lumps.ILump {
	return bsp.lumps[index]
}

func (bsp Bsp) SetLump(index int, lump lumps.ILump) Bsp {
	bsp.lumps[index] = lump
	return bsp
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
	bsp.header = readHeader(reader, bsp.header)

	//Create lumps from header data
	bsp.lumps = readLumps(reader, bsp.header, bsp.lumps)

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
	if err != nil {
		log.Fatal(err)
	}

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

		lumpData[index] = GetLumpForIndex(index, header.Version).FromBytes(raw, lumpHeader.Length)
		// Why is this here?
		// For reasons (unknown), exported data length differs from imported.
		// HOWEVER, the below snippet proves that all lumps import to export bytes are the same
		// thus ensuring validity of the process.
		/*result := lumpData[index].ToBytes()
		fmt.Println(index, len(raw), len(result))
		for i := range raw {
			if raw[i] != result[i] {
				fmt.Println(i, raw[i], result[i])
			}
		}*/
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

	for index,lump := range bsp.lumps {
		// We have to handle lump 35 (GameData differently)
		// Because valve mis-designed the file format and relatively positioned data contains absolute file offsets.
		if index == 35 {
			gamelump := bsp.lumps[index].(lumps.Game)
			gamelump = gamelump.UpdateInternalOffsets(int32(currentOffset) - bsp.header.Lumps[index].Offset)
			lumpBytes[index] = gamelump.ToBytes()
		} else {
			lumpBytes[index] = lump.ToBytes()
		}

		lumpSize := len(lumpBytes[index])

		bsp.header.Lumps[index].Length = int32(lumpSize)
		bsp.header.Lumps[index].Offset = int32(currentOffset)

		currentOffset += lumpSize
	}


	// Now we can export our bsp
	var buf bytes.Buffer

	//Write Header
	binary.Write(&buf, binary.LittleEndian, bsp.header)

	//Write lumps
	for _,lumpData := range lumpBytes {
		binary.Write(&buf, binary.LittleEndian, lumpData)

	}

	return buf.Bytes()
}

/**
	Return an instance of a Lump for a given offset.
 */
func GetLumpForIndex(index int, version int32) lumps.ILump {
	if index < 0 || index > 63 {
		log.Fatal("Invalid lump index provided.")
	}

	switch version {
	case 20:
		return versions.GetVersion20Mapping()[index]
	default:
		log.Fatal("Bsp version not currently supported")
	}
	return nil
}


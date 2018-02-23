package bsp


import (
	"bytes"
	"encoding/binary"
	"github.com/galaco/bsp/lumps"
)

// Bsp export writer.
type Writer struct {
	data Bsp
}

// Get bsp file to write.
func (w *Writer) GetBsp() Bsp {
	return w.data
}

// Set bsp file to write.
func (w *Writer) SetBsp(file Bsp) {
	w.data = file
}

// Write bsp to []byte.
func (w *Writer) Write() []byte {
	// First we need to update the header to reflect any lump changes
	// At the same time we can dump our lumps as bytes to write later
	lumpBytes := make([][]byte, 64)
	currentOffset := 1036 // Header always 1036bytes

	for index,lump := range w.data.lumps {
		// We have to handle lump 35 (GameData differently)
		// Because valve mis-designed the file format and relatively positioned data contains absolute file offsets.
		if index == 35 {
			gamelump := (lump.GetContents()).(lumps.Game)
			w.data.lumps[index].SetContents(
				gamelump.UpdateInternalOffsets(int32(currentOffset) - w.data.header.Lumps[index].Offset))
		}
		lumpBytes[index] = w.WriteLump(index)

		lumpSize := len(lumpBytes[index])

		w.data.header.Lumps[index].Length = int32(lumpSize)
		w.data.header.Lumps[index].Offset = int32(currentOffset)

		currentOffset += lumpSize

		// Finally 4byte align each lump.
		lumpBytes[index] = append(lumpBytes[index], make([]byte, currentOffset % 4)...)
		currentOffset += currentOffset % 4
	}

	// Now we can export our bsp
	var buf bytes.Buffer

	//Write Header
	binary.Write(&buf, binary.LittleEndian, w.data.header)

	//Write lumps
	for _,lumpData := range lumpBytes {
		binary.Write(&buf, binary.LittleEndian, lumpData)

	}

	return buf.Bytes()
}

/**
	Export a single lump to []byte.
 */
func (w *Writer) WriteLump(index int) []byte {
	lump := w.data.GetLump(index)
	return lump.GetContents().ToBytes()
}

/**
	Return a new bsp writer instance.
 */
func NewWriter() Writer {
	w := Writer{}
	return w
}

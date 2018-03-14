package bsp

import (
	"testing"
	"os"
	"log"
	"bytes"
)

// Test that resultant lump data matches expected.
func TestLumpExports(t *testing.T) {
	f := GetTestFile()
	length := GetBufferForTestFile(f)

	//Load file
	fileData := make([]byte, length)
	f.Read(fileData)
	f.Close()

	// Read bsp
	reader := NewReader()
	reader.SetBuffer(fileData)
	file := reader.Read()

	// Verify lump lengths
	lumpIndex := 0
	for lumpIndex < 64 {
		lump := file.GetLump(lumpIndex)
		lumpBytes := lump.GetContents().ToBytes()
		if len(lumpBytes) != int(file.GetHeader().Lumps[lumpIndex].Length) {
			t.Errorf("Lump: %d length mismatch. Got: %dbytes, expected: %dbytes", lumpIndex, len(lumpBytes), file.header.Lumps[lumpIndex].Length)
		} else {
			//fmt.Printf("Index: %d, Expected: %d, Actual: %d\n", lumpIndex, len(lump.GetRawContents()), len(lumpBytes))
			if !bytes.Equal(lumpBytes, lump.GetRawContents()) {
				t.Errorf("Lump: %d data mismatch", lumpIndex)
			}
		}

		lumpIndex += 1
	}
}

// Obtain a test file to run against.
func GetTestFile() *os.File {
	f, err := os.Open("maps/de_dust2.bsp")
	if err!= nil {
		log.Fatal(err)
	}

	return f
}

// Get a buffer to match the test file.
func GetBufferForTestFile(file *os.File) int64 {
	file.Seek(0,0)

	fi,err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	return fi.Size()
}
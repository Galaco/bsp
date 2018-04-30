package bsp

import (
	"testing"
	"os"
	"log"
	"bytes"
	"fmt"
)

// Test that resultant lump data matches expected.
func TestLumpExports(t *testing.T) {
	f := GetTestFile()


	// Read bsp
	reader := NewReader(f)
	file := reader.Read()

	//Load file
	f.Close()

	// Verify lump lengths
	lumpIndex := 0
	for lumpIndex < 64 {
		lump := file.GetLump(lumpIndex)
		lumpBytes := lump.GetContents().ToBytes()
		if len(lumpBytes) != int(file.GetHeader().Lumps[lumpIndex].Length) {
			t.Errorf("Lump: %d length mismatch. Got: %dbytes, expected: %dbytes", lumpIndex, len(lumpBytes), file.header.Lumps[lumpIndex].Length)
		} else {
			fmt.Printf("Index: %d, Expected: %d, Actual: %d\n", lumpIndex, len(lump.GetRawContents()), len(lumpBytes))
			if !bytes.Equal(lumpBytes, lump.GetRawContents()) {
				t.Errorf("Lump: %d data mismatch", lumpIndex)
			}
		}

		lumpIndex += 1
	}
}

// Obtain a test file to run against.
func GetTestFile() *os.File {
	f, err := os.Open("maps/ze_bioshock_v6_2.bsp")
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
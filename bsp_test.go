package bsp

import (
	"testing"
	"log"
	"bytes"
)

// Test that resultant lump data matches expected.
func TestLumpExports(t *testing.T) {
	file,err := ReadFromFile("maps/v20/de_dust2.bsp")
	if err != nil {
		t.Error(err)
	}

	// Verify lump lengths
	lumpIndex := 0
	for lumpIndex < 64 {
		lump := file.GetLump(lumpIndex)
		lumpBytes := lump.GetContents().ToBytes()
		if len(lumpBytes) != int(file.GetHeader().Lumps[lumpIndex].Length) {
			t.Errorf("Lump: %d length mismatch. Got: %dbytes, expected: %dbytes", lumpIndex, len(lumpBytes), file.header.Lumps[lumpIndex].Length)
		} else {
			log.Printf("Index: %d, Expected: %d, Actual: %d\n", lumpIndex, len(lump.GetRawContents()), len(lumpBytes))
			if !bytes.Equal(lumpBytes, lump.GetRawContents()) {
				t.Errorf("Lump: %d data mismatch", lumpIndex)
			}
		}

		lumpIndex += 1
	}
}
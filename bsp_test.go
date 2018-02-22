package bsp

import (
	"testing"
	"os"
	"log"
)

/**
	Test that resultant lump data matches expected.
 */
func TestLumpExports(t *testing.T) {
	f := GetTestFile()
	file := Parse(f)
	f.Close()

	// Verify lump lengths
	lumpIndex := 0
	for lumpIndex < 64 {
		lump := file.GetLump(lumpIndex)
		lumpBytes := lump.ToBytes()
		if len(lumpBytes) != int(file.header.Lumps[lumpIndex].Length) {
			t.Errorf("Lump: %d length mismatch. Got: %dbytes, expected: %dbytes", lumpIndex, len(lumpBytes), file.header.Lumps[lumpIndex].Length)
		}
		lumpIndex += 1
	}
}




func GetTestFile() *os.File {
	f, err := os.Open("ze_bioshock_v6_2.bsp")
	if err!= nil {
		log.Fatal(err)
	}

	return f
}
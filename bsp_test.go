package bsp

import (
	"bytes"
	"log"
	"testing"
)

// Test that resultant lump data matches expected.
func TestLumpExports(t *testing.T) {
	t.Skip()
	file, err := ReadFromFile("maps/v20/de_dust2.bsp")
	if err != nil {
		t.Error(err)
	}

	// Verify lump lengths
	lumpIndex := 0
	for lumpIndex < 64 {
		lump := file.GetLump(LumpId(lumpIndex))
		rawLump := file.GetLumpRaw(LumpId(lumpIndex))
		lumpBytes, err := lump.Marshall()
		if err != nil {
			t.Error(err)
		}
		if len(lumpBytes) != int(file.GetHeader().Lumps[lumpIndex].Length) {
			t.Errorf("Lump: %d length mismatch. Got: %dbytes, expected: %dbytes", lumpIndex, len(lumpBytes), file.header.Lumps[lumpIndex].Length)
		} else {
			log.Printf("Index: %d, Expected: %d, Actual: %d\n", lumpIndex, len(rawLump.GetRawContents()), len(lumpBytes))
			if !bytes.Equal(lumpBytes, rawLump.GetRawContents()) {
				t.Errorf("Lump: %d data mismatch", lumpIndex)
			}
		}

		lumpIndex += 1
	}
}

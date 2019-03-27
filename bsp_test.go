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
		lump := file.Lump(LumpId(lumpIndex))
		rawLump := file.RawLump(LumpId(lumpIndex))
		lumpBytes, err := lump.Marshall()
		if err != nil {
			t.Error(err)
		}
		if len(lumpBytes) != int(file.Header().Lumps[lumpIndex].Length) {
			t.Errorf("Lump: %d length mismatch. Got: %dbytes, expected: %dbytes", lumpIndex, len(lumpBytes), file.header.Lumps[lumpIndex].Length)
		} else {
			log.Printf("Index: %d, Expected: %d, Actual: %d\n", lumpIndex, len(rawLump.RawContents()), len(lumpBytes))
			if !bytes.Equal(lumpBytes, rawLump.RawContents()) {
				t.Errorf("Lump: %d data mismatch", lumpIndex)
			}
		}

		lumpIndex += 1
	}
}

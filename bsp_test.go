package bsp

import (
	"testing"
	"os"
	"log"
)

func TestToBytes(t *testing.T) {
	f := GetTestFile()

	info,err := f.Stat()
	if err!= nil {
		log.Fatal(err)
	}
	inFileBytes := make([]byte, info.Size())
	f.Read(inFileBytes)


	file := Parse(f)
	f.Close()

	outFileBytes := ToBytes(file)

	if len(inFileBytes) != len(outFileBytes) {
		t.Errorf("Export length mismatch. Got: %dbytes, expected: %dbytes", len(outFileBytes), len(inFileBytes))
	}
}

func TestLumpIntegrity(t *testing.T) {
	f := GetTestFile()
	file := Parse(f)
	f.Close()

	lumpIndex := 0
	for lumpIndex < 64 {
		lump := file.GetLump(lumpIndex)
		lumpBytes := lump.ToBytes()
		if len(lumpBytes) != int(file.header.Lumps[lumpIndex].Length) {
			t.Errorf("Lump %d length mismatch. Got: %dbytes, expected: %dbytes", len(lumpBytes), file.header.Lumps[lumpIndex].Length)
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
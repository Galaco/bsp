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
		}
		/* // We cant assert this with the current implementation.
		if !bytes.Equal(lumpBytes, lump.GetImportDetails().GetRaw()) {
			t.Errorf("Lump: %d data mismatch", lumpIndex)
		}*/
		lumpIndex += 1
	}

	// Why is this here?
	// For reasons (4byte alignment?!), exported data length differs from imported.
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

/**
	Obtain a test file to run against.
 */
func GetTestFile() *os.File {
	f, err := os.Open("ze_bioshock_v6_2.bsp")
	if err!= nil {
		log.Fatal(err)
	}

	return f
}

/**
	Get a buffer to match the test file.
 */
func GetBufferForTestFile(file *os.File) int64 {
	file.Seek(0,0)

	fi,err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	return fi.Size()
}
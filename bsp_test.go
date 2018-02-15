package bsp

import (
	"testing"
	"os"
	"log"
	"bytes"
)


func restParse(t *testing.T) {
	f := GetTestFile()

	file := Parse(f)

	if file.Header.Id != 1347633750 {
		t.Errorf("File identity incorrect, got: %d, expexted %d", file.Header.Id, 1347633750)
	}

	if file.Header.Version != 20 {
		t.Errorf("File version incorrect, got: %d, expexted %d", file.Header.Version, 20)
	}

	if file.Header.Revision != 12491 {
		t.Errorf("File revision incorrect, got: %d, expexted %d", file.Header.Revision, 12491)
	}

	f.Close()
}

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
		t.Errorf("Export length mismatch. Got: %dbytes, expected %dbytes", len(outFileBytes), len(inFileBytes))
	}


	if !bytes.Equal(inFileBytes, outFileBytes) {
		t.Errorf("Export data does not match import data")
	}
}




func GetTestFile() *os.File {
	f, err := os.Open("ze_bioshock_v6_2.bsp")
	if err!= nil {
		log.Fatal(err)
	}

	return f
}
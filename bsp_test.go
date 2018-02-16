package bsp

import (
	"testing"
	"os"
	"log"
	"bytes"
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
		t.Errorf("Export length mismatch. Got: %dbytes, expected %dbytes", len(outFileBytes), len(inFileBytes))
	}


	if !bytes.Equal(inFileBytes, outFileBytes) {
		t.Errorf("Export data does not match import data")
	}
}




func GetTestFile() *os.File {
	f, err := os.Open("aim_office_battle.bsp")
	if err!= nil {
		log.Fatal(err)
	}

	return f
}
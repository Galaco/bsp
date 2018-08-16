package bsp

import (
	"testing"
	"os"
)

func TestReadFromFile(t *testing.T) {
	_,err := ReadFromFile("maps/v20/de_dust2.bsp")
	if err != nil {
		t.Error(err)
	}
}

func TestReadFromStream(t *testing.T) {
	f, err := os.Open("maps/v20/de_dust2.bsp")
	if err != nil {
		t.Error(err)
	}

	_,err = ReadFromStream(f)
	if err != nil {
		t.Error(err)
	}
}

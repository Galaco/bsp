package bsp

import (
	"github.com/galaco/bsp/lumps"
	"os"
	"testing"
)

func TestReadFromFile(t *testing.T) {
	t.Skip()
	_, err := ReadFromFile("maps/v20/de_dust2.bsp")
	if err != nil {
		t.Error(err)
	}
}

func TestReadFromStream(t *testing.T) {
	t.Skip()
	f, err := os.Open("maps/v20/de_dust2.bsp")
	if err != nil {
		t.Error(err)
	}

	r, err := ReadFromStream(f)
	if err != nil {
		t.Error(err)
	}

	r.Lump(LumpGame).(*lumps.Game).GetStaticPropLump()
}

package bsp

import (
	"github.com/galaco/bsp/lumps"
	"os"
	"testing"
)

func TestReadFromFile(t *testing.T) {
	_, err := ReadFromFile("maps/v20/de_dust2.bsp")
	if err != nil {
		t.Error(err)
	}
}

func TestReadFromStream(t *testing.T) {
	f, err := os.Open("maps/v20/de_dust2.bsp")
	if err != nil {
		t.Error(err)
	}

	r, err := ReadFromStream(f)
	if err != nil {
		t.Error(err)
	}

	r.GetLump(LUMP_GAME_LUMP).(*lumps.Game).GetStaticPropLump()
}

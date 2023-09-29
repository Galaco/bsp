package main

import (
	"compress/gzip"
	"github.com/galaco/bsp"
	"github.com/galaco/bsp/lump"
	"os"
)

func main() {
	f, err := os.Open("testdata/v20/de_dust2.bsp.gz")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	binarygzr, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}

	reader := bsp.NewReader()
	file, err := reader.Read(binarygzr)
	if err != nil {
		panic(err)
	}

	// Make a modification to a lump.
	file.Lumps[bsp.LumpMapFlags].(*lump.MapFlags).Data.LevelFlags &= 8

	writer := bsp.NewWriter()
	_, err = writer.Write(file)
	if err != nil {
		panic(err)
	}
}

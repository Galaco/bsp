package main

import (
	"compress/gzip"
	"github.com/galaco/bsp"
	"github.com/galaco/bsp/lump"
	"log"
	"os"
)

func main() {
	f, err := os.Open("testdata/v20/de_dust2.bsp.gz")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	binarygzr, err := gzip.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}

	reader := bsp.NewReader()
	file, err := reader.Read(binarygzr)
	if err != nil {
		log.Fatal(err)
	}

	entdata := file.Lumps[bsp.LumpEntities].(*lump.EntData)
	log.Println(entdata.Contents())
}

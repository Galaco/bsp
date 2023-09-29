package main

import (
	"compress/gzip"
	"github.com/galaco/bsp"
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
	_ = file
}

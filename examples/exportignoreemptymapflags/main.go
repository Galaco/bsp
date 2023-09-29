package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/galaco/bsp"
	"github.com/galaco/bsp/lump"
	"log"
	"os"
)

type CustomMapFlags struct {
	lump.MapFlags
}

func (mapFlags *CustomMapFlags) ToBytes() ([]byte, error) {
	// If our flags are 0, return an empty lump.
	// CS:S de_nuke works this way, but de_dust2 does not.
	if mapFlags.MapFlags.Data.LevelFlags == 0 {
		return make([]byte, 0), nil
	}
	return mapFlags.MapFlags.ToBytes()
}

func customLumpResolver(id bsp.LumpId, header bsp.Header) (l bsp.Lump, err error) {
	if id == bsp.LumpMapFlags {
		return &CustomMapFlags{}, nil
	}
	return bsp.LumpResolverByBSPVersion(id, header)
}

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

	reader := bsp.NewReaderWithConfig(bsp.ReaderConfig{
		LumpResolver: customLumpResolver,
	})
	file, err := reader.Read(binarygzr)
	if err != nil {
		panic(err)
	}
	log.Println(fmt.Sprintf("Old MapFlags Length: %d", file.Header.Lumps[bsp.LumpMapFlags].Length))

	// Make a modification to a lump.
	file.Lumps[bsp.LumpMapFlags].(*CustomMapFlags).Data.LevelFlags = 0

	writer := bsp.NewWriter()
	output, err := writer.Write(file)
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(output)
	reader2 := bsp.NewReader()
	newFile, err := reader2.Read(buf)
	if err != nil {
		panic(err)
	}
	log.Println(fmt.Sprintf("New MapFlags Length: %d", newFile.Header.Lumps[bsp.LumpMapFlags].Length))
}

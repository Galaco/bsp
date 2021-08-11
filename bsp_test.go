package bsp

import (
	"bytes"
	"hash/crc32"
	"log"
	"sort"
	"testing"

	"github.com/galaco/bsp/lumps"
)

// Test that resultant lump data matches expected.
func TestLumpExports(t *testing.T) {
	t.Skip()
	file, err := ReadFromFile("maps/v20/de_dust2.bsp")
	if err != nil {
		t.Error(err)
	}

	// Verify lump lengths
	lumpIndex := 0
	for lumpIndex < 64 {
		lump := file.Lump(LumpId(lumpIndex))
		rawLump := file.RawLump(LumpId(lumpIndex))
		lumpBytes, err := lump.Marshall()
		if err != nil {
			t.Error(err)
		}
		if len(lumpBytes) != int(file.Header().Lumps[lumpIndex].Length) {
			t.Errorf("Lump: %d length mismatch. Got: %dbytes, expected: %dbytes", lumpIndex, len(lumpBytes), file.header.Lumps[lumpIndex].Length)
		} else {
			log.Printf("Index: %d, Expected: %d, Actual: %d\n", lumpIndex, len(rawLump.RawContents()), len(lumpBytes))
			if !bytes.Equal(lumpBytes, rawLump.RawContents()) {
				t.Errorf("Lump: %d data mismatch", lumpIndex)
			}
		}

		lumpIndex += 1
	}
}

func TestBsp_Header(t *testing.T) {
	sut := new(Bsp)
	header := new(Header)
	header.Id = 564

	sut.header = *header

	if sut.Header().Id != header.Id {
		t.Error("unexpected header struct returned")
	}
}

func TestBsp_Lump(t *testing.T) {
	sut := new(Bsp)
	l := new(lumps.Generic)
	sut.lumps[0].data = l

	if sut.Lump(0) != l {
		t.Error("unexpected lump returned")
	}
}

func TestBsp_RawLump(t *testing.T) {
	sut := new(Bsp)
	l := new(lumps.Generic)
	sut.lumps[0].data = l

	if sut.RawLump(0).data != l {
		t.Error("unexpected lump returned")
	}
}

func TestBsp_SetLump(t *testing.T) {
	sut := new(Bsp)
	ld := new(lumps.Generic)
	l := new(Lump)
	l.data = ld
	sut.SetLump(0, *l)

	if sut.RawLump(0).data != ld {
		t.Error("unexpected lump returned")
	}
}

func TestBsp_Crc(t *testing.T) {
	file, err := ReadFromFile("de_inferno.bsp")
	if err != nil {
		t.Error(err)
	}

	const lumpCount = 64

	crc := crc32.NewIEEE()

	lumpList := make([]LumpId, lumpCount)
	for i := 0; i < lumpCount; i++ {
		lumpList[i] = LumpId(i)
	}

	sort.Slice(lumpList, func(i, j int) bool {
		return file.header.Lumps[i].Offset < file.header.Lumps[j].Offset
	})

	for i := 0; i < lumpCount; i++ {
		l := lumpList[i]
		if l == LumpEntities {
			continue
		}

		b := file.RawLump(l).raw

		if len(b) != int(file.header.Lumps[l].Length) {
			panic("???")
		}

		_, err := crc.Write(b)
		if err != nil {
			panic(err)
		}
	}

	if crc.Sum32() != uint32(4271092677) {
		panic(crc.Sum32())
	}
}

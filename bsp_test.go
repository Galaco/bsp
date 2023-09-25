package bsp

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"os"
	"testing"

	"github.com/galaco/bsp/lumps"
	"github.com/google/go-cmp/cmp"
)

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
	l := new(lumps.RawBytes)
	sut.lumps[0] = l

	if sut.Lump(0) != l {
		t.Error("unexpected lump returned")
	}
}

// Test that resultant lump data matches expected.
func TestLumpExports(t *testing.T) {
	testCases := []struct {
		name     string
		filePath string
	}{
		//{
		//	name:     "de_dust2",
		//	filePath: "testdata/v20/de_dust2.bsp.gz",
		//},
		{
			name:     "ar_baggage",
			filePath: "testdata/v21/ar_baggage.bsp.gz",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			binaryFile, err := os.Open(tc.filePath)
			if err != nil {
				t.Error(err)
			}
			binarygzr, err := gzip.NewReader(binaryFile)
			if err != nil {
				t.Error(err)
			}
			binaryData, err := io.ReadAll(binarygzr)
			if err != nil {
				t.Error(err)
			}

			testFile, err := os.Open(tc.filePath)
			if err != nil {
				t.Error(err)
			}
			testFilegzr, err := gzip.NewReader(testFile)
			if err != nil {
				t.Error(err)
			}
			testBSP, err := ReadFromStream(testFilegzr)
			if err != nil {
				t.Error(err)
			}

			// Verify lump lengths.
			lumpIndex := 0
			for lumpIndex < 64 {
				actual, err := testBSP.Lump(LumpId(lumpIndex)).ToBytes()
				if err != nil {
					t.Error(err)
				}
				if len(actual) != int(testBSP.Header().Lumps[lumpIndex].Length) {
					t.Errorf("Lump: %d length mismatch. Got: %dbytes, expected: %dbytes", lumpIndex, len(actual), testBSP.header.Lumps[lumpIndex].Length)
				}

				expected := binaryData[testBSP.Header().Lumps[lumpIndex].Offset : testBSP.Header().Lumps[lumpIndex].Offset+testBSP.Header().Lumps[lumpIndex].Length]
				if !bytes.Equal(actual, expected) {
					t.Errorf("Lump: %d data mismatch", lumpIndex)
					log.Println(cmp.Diff(expected, actual))
				}

				lumpIndex += 1
			}
		})
	}
}

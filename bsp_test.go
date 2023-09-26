package bsp

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Test that resultant lump data matches expected.
func Test_ExportedLumpBytesAreCorrect(t *testing.T) {
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
			defer binaryFile.Close()
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
			defer testFile.Close()
			testFilegzr, err := gzip.NewReader(testFile)
			if err != nil {
				t.Error(err)
			}
			testBSP, err := NewReader().Read(testFilegzr)
			if err != nil {
				t.Error(err)
			}

			// Verify lump lengths.
			for lumpIndex := 0; lumpIndex < 64; lumpIndex++ {
				actual, err := testBSP.Lumps[LumpId(lumpIndex)].ToBytes()
				if err != nil {
					t.Error(err)
				}
				if len(actual) != int(testBSP.Header.Lumps[lumpIndex].Length) {
					t.Errorf("Lump: %d length mismatch. Got: %dbytes, expected: %dbytes", lumpIndex, len(actual), testBSP.Header.Lumps[lumpIndex].Length)
				}

				expected := binaryData[testBSP.Header.Lumps[lumpIndex].Offset : testBSP.Header.Lumps[lumpIndex].Offset+testBSP.Header.Lumps[lumpIndex].Length]
				if !bytes.Equal(actual, expected) {
					t.Errorf("Lump: %d data mismatch", lumpIndex)
					log.Println(cmp.Diff(expected, actual))
				}
			}
		})
	}
}

package bsp

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io"
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewWriter(t *testing.T) {
	w := NewWriter()
	if w == nil {
		t.Error("NewWriter returned nil")
	}
}

func TestWriter_Write(t *testing.T) {
	testCases := []struct {
		name     string
		filePath string
	}{
		{
			name:     "de_dust2",
			filePath: "testdata/v20/de_dust2.bsp.gz",
		},
		{
			name:     "ar_baggage",
			filePath: "testdata/v21/ar_baggage.bsp.gz",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.filePath)
			if err != nil {
				t.Error(err)
			}
			defer func(f *os.File) {
				if err := f.Close(); err != nil {
					t.Error(err)
				}
			}(f)
			binarygzr, err := gzip.NewReader(f)
			if err != nil {
				t.Error(err)
			}

			var expected bytes.Buffer
			bsp, err := NewReader(LumpResolverByBSPVersion()).Read(io.TeeReader(binarygzr, &expected))
			if err != nil {
				t.Error(err)
			}
			expectedBytes := expected.Bytes()

			actual, err := NewWriter().toBytes(bsp)
			if err != nil {
				t.Fatalf("toBytes(%s) returned error: %v", tc.filePath, err)
			}

			if !bytes.Equal(expectedBytes, actual) {
				t.Errorf("toBytes(%s) returned unexpected bytes", tc.filePath)

				// It's rather painful to easily debug this kind of output,
				// below are some helpers for diffing header and lump data.
				baseOffset := 8
				compSize := 16
				for _, i := range resolveLumpExportOrder(&bsp.Header) {
					offset := baseOffset + (int(i) * compSize)

					expectedLumpOffset := int(binary.LittleEndian.Uint32(expectedBytes[offset : offset+4]))
					expectedLumpLength := int(binary.LittleEndian.Uint32(expectedBytes[offset+4 : offset+8]))
					actualLumpOffset := int(binary.LittleEndian.Uint32(expectedBytes[offset : offset+4]))
					actualLumpLength := int(binary.LittleEndian.Uint32(expectedBytes[offset+4 : offset+8]))

					log.Printf("%d: offset: %d, %d (%d). length: %d, %d (%d). Version: %d, %d\n", i,
						expectedLumpOffset,
						actualLumpOffset,
						actualLumpOffset-expectedLumpOffset,
						expectedLumpLength,
						actualLumpLength,
						expectedLumpLength-actualLumpLength,
						binary.LittleEndian.Uint32(expectedBytes[offset+8:offset+12]),
						binary.LittleEndian.Uint32(actual[offset+8:offset+12]),
					)

					// Compare header
					if !bytes.Equal(expectedBytes[offset:offset+compSize], actual[offset:offset+compSize]) {
						t.Errorf("%d: header.toBytes(%s) returned unexpected bytes", i, tc.filePath)

						if diff := cmp.Diff(expectedBytes[offset:offset+compSize], actual[offset:offset+compSize]); diff != "" {
							t.Errorf("%d: header.toBytes(%s) returned unexpected diff (-want +got):\n%s", i, tc.filePath, diff)
						}
					}

					// Compare lump the header references.
					if !bytes.Equal(expectedBytes[expectedLumpOffset:expectedLumpOffset+expectedLumpLength], actual[actualLumpOffset:actualLumpOffset+actualLumpLength]) {
						t.Errorf("%d: lump.toBytes(%s) returned unexpected bytes", i, tc.filePath)

						if diff := cmp.Diff(expectedBytes[expectedLumpOffset:expectedLumpOffset+expectedLumpLength], actual[actualLumpOffset:actualLumpOffset+actualLumpLength]); diff != "" {
							t.Errorf("%d: header.toBytes(%s) returned unexpected diff (-want +got):\n%s", i, tc.filePath, diff)
						}
					}
				}

				// now compare each body.
				if len(expectedBytes) != len(actual) {
					t.Errorf("toBytes(%s) returned unexpected bytes", tc.filePath)
				}

				//if diff := cmp.Diff(expectedBytes, actual); diff != "" {
				//	t.Errorf("toBytes(%s) returned unexpected diff (-want +got):\n%s", tc.filePath, diff)
				//}
			}
		})
	}
}

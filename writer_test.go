package bsp

import (
	"bytes"
	"compress/gzip"
	"io"
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
			bsp, err := NewReader().Read(io.TeeReader(binarygzr, &expected))
			if err != nil {
				t.Error(err)
			}
			expectedBytes := expected.Bytes()

			actual, err := NewWriter().toBytes(bsp)
			if err != nil {
				t.Fatalf("toBytes(%s) returned error: %v", tc.filePath, err)
			}

			//baseOffset := 8 - 1
			//compSize := 16
			//for i := 0; i < 64; i++ {
			//	expect := binary.BigEndian.Uint32(expectedBytes[baseOffset+(i*compSize) : baseOffset+(i*compSize)+4])
			//	act := binary.BigEndian.Uint32(actual[baseOffset+(i*compSize) : baseOffset+(i*compSize)+4])
			//	log.Printf("%d: %d, %d\n", i, expect, act)
			//	if !bytes.Equal(expectedBytes[baseOffset+(i*compSize):baseOffset+(i*(compSize*2))], actual[baseOffset+(i*compSize):baseOffset+(i*(compSize*2))]) {
			//		t.Errorf("%d: toBytes(%s) returned unexpected bytes", i, tc.filePath)
			//	}
			//
			//	if diff := cmp.Diff(expectedBytes[7+(i*16):7+(i*32)], actual[7+(i*16):7+(i*32)]); diff != "" {
			//		t.Errorf("%d: toBytes(%s) returned unexpected diff (-want +got):\n%s", i, tc.filePath, diff)
			//	}
			//}

			if !bytes.Equal(expectedBytes[8:1024], actual[8:1024]) {
				t.Errorf("toBytes(%s) returned unexpected bytes", tc.filePath)
			}
			if diff := cmp.Diff(expectedBytes[8:1024], actual[8:1024]); diff != "" {
				t.Errorf("toBytes(%s) returned unexpected diff (-want +got):\n%s", tc.filePath, diff)
			}
		})
	}
}

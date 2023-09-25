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
			defer f.Close()
			binarygzr, err := gzip.NewReader(f)
			if err != nil {
				t.Error(err)
			}

			var buf bytes.Buffer
			bsp, err := ReadFromStream(io.TeeReader(binarygzr, &buf))
			if err != nil {
				t.Error(err)
			}

			w := NewWriter()
			actual, err := w.Write(bsp)
			if err != nil {
				t.Fatalf("Write(%s) returned error: %v", tc.filePath, err)
			}
			expected := buf.Bytes()
			if !bytes.Equal(expected, actual) {
				t.Errorf("Write(%s) returned unexpected bytes", tc.filePath)
			}
			if diff := cmp.Diff(expected[:256], actual[:256]); diff != "" {
				t.Errorf("Write(%s) returned unexpected diff (-want +got):\n%s", tc.filePath, diff)
			}
		})
	}
}

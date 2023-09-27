package bsp_test

import (
	"compress/gzip"
	"os"
	"testing"

	"github.com/galaco/bsp"
)

func TestBsp_Crc(t *testing.T) {
	testCases := []struct {
		name     string
		filePath string
		expected uint32
	}{
		{
			name:     "de_dust2",
			filePath: "testdata/v20/de_dust2.bsp.gz",
			expected: 1947027563,
		},
		{
			name:     "ar_baggage",
			filePath: "testdata/v21/ar_baggage.bsp.gz",
			expected: 2836609078,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.filePath)
			if err != nil {
				t.Error(err)
			}

			gzR, err := gzip.NewReader(f)
			if err != nil {
				t.Error(err)
			}

			bspF, err := bsp.NewReader().Read(gzR)
			if err != nil {
				t.Error(err)
			}

			res, err := bspF.CRC32()
			if err != nil {
				t.Error("unexpected error:", err)
			}

			if res != tc.expected {
				t.Errorf("CRC incorrect, expected %d got %d", tc.expected, res)
			}
		})
	}
}

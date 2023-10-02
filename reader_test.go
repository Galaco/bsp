package bsp

import (
	"compress/gzip"
	"os"
	"testing"

	"github.com/galaco/bsp/lump"
)

func TestReadFromStream(t *testing.T) {
	testCases := []struct {
		name          string
		filePath      string
		expectedError error
	}{
		{
			name:          "de_dust2",
			filePath:      "testdata/v20/de_dust2.bsp.gz",
			expectedError: nil,
		},
		{
			name:          "ar_baggage",
			filePath:      "testdata/v21/ar_baggage.bsp.gz",
			expectedError: nil,
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

			r, err := NewReaderWithConfig(ReaderConfig{
				LumpResolver: LumpResolverByBSPVersion,
			}).Read(binarygzr)
			if err != nil {
				t.Error(err)
			}

			r.Lumps[LumpGame].(*lump.Game).GetStaticPropLump()
		})
	}
}

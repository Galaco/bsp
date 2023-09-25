package bsp

import (
	"os"
	"testing"

	"github.com/galaco/bsp/lumps"
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.filePath)
			if err != nil {
				t.Error(err)
			}

			r, err := ReadFromStream(f)
			if err != nil {
				t.Error(err)
			}

			r.Lump(LumpGame).(*lumps.Game).GetStaticPropLump()
		})
	}
}

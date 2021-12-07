package bsp_test

import (
	"compress/gzip"
	"os"
	"testing"

	"github.com/galaco/bsp"
)

func TestBsp_Crc(t *testing.T) {
	f, err := os.Open("ar_baggage.bsp.gz")
	if err != nil {
		t.Error(err)
	}

	gzR, err := gzip.NewReader(f)
	if err != nil {
		t.Error(err)
	}

	bspF, err := bsp.ReadFromStream(gzR)
	if err != nil {
		t.Error(err)
	}

	res, err := bspF.CRC32()
	if err != nil {
		t.Error("unexpected error:", err)
	}

	const expected = 2836609078

	if res != expected {
		t.Errorf("CRC incorrect, expected %d got %d", expected, res)
	}
}

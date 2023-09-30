package bsp

import (
	"reflect"
	"testing"

	"github.com/galaco/bsp/lump"
)

func TestGetLumpForVersion(t *testing.T) {
	testCases := []struct {
		name    string
		id      int
		version int
		t       reflect.Type
	}{
		{
			name:    "unknown version",
			id:      4,
			version: 987,
			t:       reflect.TypeOf(&lump.RawBytes{}),
		},
		{
			name:    "v19",
			id:      4,
			version: 19,
			t:       reflect.TypeOf(&lump.Visibility{}),
		},
		{
			name:    "v20",
			id:      4,
			version: 20,
			t:       reflect.TypeOf(&lump.Visibility{}),
		},
		{
			name:    "v21",
			id:      4,
			version: 21,
			t:       reflect.TypeOf(&lump.Visibility{}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l, err := LumpResolverByBSPVersion(LumpId(tc.id), Header{
				Version: int32(tc.version),
			})
			if err != nil {
				t.Error(err)
			}
			if reflect.TypeOf(l) != tc.t {
				t.Error("Lump type mismatch.")
			}
		})
	}
}

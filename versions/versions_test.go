package versions

import (
	"reflect"
	"testing"

	"github.com/galaco/bsp/lumps"
)

func TestGetLumpForVersion(t *testing.T) {
	l, err := GetLumpForVersion(987, 1)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(l) != reflect.TypeOf(&lumps.RawBytes{}) {
		t.Error("Lump type mismatch.")
	}

	l, err = GetLumpForVersion(19, 4)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(l) != reflect.TypeOf(&lumps.Visibility{}) {
		t.Error("Lump type mismatch.")
	}

	l, err = GetLumpForVersion(20, 4)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(l) != reflect.TypeOf(&lumps.Visibility{}) {
		t.Error("Lump type mismatch.")
	}

	l, err = GetLumpForVersion(21, 4)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(l) != reflect.TypeOf(&lumps.Visibility{}) {
		t.Error("Lump type mismatch.")
	}
}

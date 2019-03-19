package versions

import (
	"github.com/galaco/bsp/lumps"
	"reflect"
	"testing"
)

func TestGetLumpForVersion(t *testing.T) {
	l,err := GetLumpForVersion(987, 1)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(l) != reflect.TypeOf(&lumps.Unimplemented{}) {
		t.Error("Lump type mismatch.")
	}

	l,err = GetLumpForVersion(19, 4)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(l) != reflect.TypeOf(&lumps.Visibility{}) {
		t.Error("Lump type mismatch.")
	}

	l,err = GetLumpForVersion(20, 4)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(l) != reflect.TypeOf(&lumps.Visibility{}) {
		t.Error("Lump type mismatch.")
	}

	l,err = GetLumpForVersion(21, 4)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(l) != reflect.TypeOf(&lumps.Visibility{}) {
		t.Error("Lump type mismatch.")
	}
}

package versions

import (
	"reflect"
	"testing"

	"github.com/galaco/bsp/lumps"
)

func TestGetv20Lump(t *testing.T) {
	l, _ := Getv20Lump(4)
	if reflect.TypeOf(l) != reflect.TypeOf(&lumps.Visibility{}) {
		t.Errorf("Lump type mismatch. Got: %s, expected: %s, ",
			reflect.TypeOf(l),
			reflect.TypeOf(lumps.Visibility{}))
	}

	_, err := Getv20Lump(65)
	if err == nil {
		t.Error("invalid lump index provided, but not error returned")
	}
}

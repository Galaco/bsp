package versions

import (
	"reflect"
	"testing"

	"github.com/galaco/bsp/lumps"
)

func TestGetv19Lump(t *testing.T) {
	l, _ := Getv19Lump(4)
	if reflect.TypeOf(l) != reflect.TypeOf(&lumps.Visibility{}) {
		t.Errorf("Lump type mismatch. Got: %s, expected: %s, ",
			reflect.TypeOf(l),
			reflect.TypeOf(lumps.Visibility{}))
	}
}

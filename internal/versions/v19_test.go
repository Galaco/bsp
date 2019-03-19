package versions

import (
	"github.com/galaco/bsp/lumps"
	"reflect"
	"testing"
)

func TestGetv19Lump(t *testing.T) {
	l, _ := Getv19Lump(4)
	if reflect.TypeOf(l) != reflect.TypeOf(&lumps.Visibility{}) {
		t.Errorf("Lump type mismatch. Got: %s, expected: %s, ",
			reflect.TypeOf(l),
			reflect.TypeOf(lumps.Visibility{}))
	}
}

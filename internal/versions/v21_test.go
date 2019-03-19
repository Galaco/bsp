package versions

import (
	"github.com/galaco/bsp/lumps"
	"reflect"
	"testing"
)

func TestGetv21Lump(t *testing.T) {
	l, _ := Getv21Lump(4)
	if reflect.TypeOf(l) != reflect.TypeOf(&lumps.Visibility{}) {
		t.Errorf("Lump type mismatch. Got: %s, expected: %s, ",
			reflect.TypeOf(l),
			reflect.TypeOf(lumps.Visibility{}))
	}
}

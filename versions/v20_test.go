package versions

import (
	"github.com/galaco/bsp/lumps"
	"reflect"
	"testing"
)

func TestLumpIndexReturnType(t *testing.T) {
	l, _ := Getv20Lump(4)
	if reflect.TypeOf(l) != reflect.TypeOf(&lumps.Visibility{}) {
		t.Errorf("Lump type mismatch. Got: %s, expected: %s, ",
			reflect.TypeOf(l),
			reflect.TypeOf(lumps.Visibility{}))
	}
}

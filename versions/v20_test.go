package versions

import (
	"testing"
	"github.com/galaco/bsp/lumps"
	"reflect"
)

func TestLumpIndexReturnType(t *testing.T) {
	if reflect.TypeOf(GetVersion20Mapping()[4]) != reflect.TypeOf(lumps.Visibility{}) {
		t.Errorf("Lump type mismatch. Got: %s, expected: %s, ",
			reflect.TypeOf(GetVersion20Mapping()[4]),
			reflect.TypeOf(lumps.Visibility{}))
	}
}
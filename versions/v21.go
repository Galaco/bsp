package versions

import (
	"github.com/galaco/bsp/lumps"
)

func Getv21Lump(index int) (lumps.ILump,error) {
	var ret lumps.ILump
	var err error
	switch index {
	default:
		return Getv20Lump(index)
	}

	return ret,err
}

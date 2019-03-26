package versions

import (
	"github.com/galaco/bsp/lumps"
)

// Getv21Lump returns the corresponding v21 lump for provided index
func Getv21Lump(index int) (lumps.ILump, error) {
	return Getv20Lump(index)
}

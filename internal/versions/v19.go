package versions

import (
	"github.com/galaco/bsp/lumps"
)

// Getv19Lump returns the corresponding v19 lump for provided index
func Getv19Lump(index int) (lumps.ILump, error) {
	return Getv20Lump(index)
}

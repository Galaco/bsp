package versions

import (
	"github.com/galaco/bsp/lumps"
)

// Returns an empty bsp v19 lump for the specified lump id
func Getv21Lump(index int) (lumps.ILump, error) {
	return Getv20Lump(index)
}

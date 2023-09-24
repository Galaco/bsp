package versions

import "github.com/galaco/bsp/lumps"

// GetLumpForVersion returns an empty bsp lump for the specified bsp version and lump id
// If a version is not 19,20,21 then a generic lump that holds
// raw bytes only ([]byte) is returned.
func GetLumpForVersion(bspVersion int, lumpId int) (lumps.ILump, error) {
	switch bspVersion {
	case 19:
		return Getv19Lump(lumpId)
	case 20:
		return Getv20Lump(lumpId)
	case 21:
		return Getv21Lump(lumpId)
	default:
		return new(lumps.RawBytes), nil
	}
}

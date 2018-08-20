package versions

import "github.com/galaco/bsp/lumps"

func GetLumpForVersion(bspVersion int, lumpId int) (lumps.ILump,error) {
	switch bspVersion {
	case 20:
		return Getv20Lump(lumpId)
	case 21:
		return Getv21Lump(lumpId)
	default:
		return &lumps.Unimplemented{},nil
	}
}

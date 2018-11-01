package game

const StaticPropLumpId = 1936749168

type Header struct {
	LumpCount int32
	GameLumps []LumpDef // Slice length must equal lumpCount. Validation to be added
}

func (header Header) SetLumpCount(num int32) Header {
	header.LumpCount = num
	header.GameLumps = make([]LumpDef, header.LumpCount)
	return header
}

type LumpDef struct {
	Id         int32
	Flags      uint16
	Version    uint16
	FileOffset int32
	FileLength int32
}

type GenericGameLump struct {
	Length int32
	Data   []byte
}

package lumpdata


type GameHeader struct {
	LumpCount int32
	GameLumps []GameLump // Slice length must equal lumpCount. Validation to be added
}
func (header GameHeader) SetLumpCount(num int32) GameHeader{
	header.LumpCount = num
	header.GameLumps = make([]GameLump, header.LumpCount)
	return header
}

type GameLump struct {
	Id int32
	Flags uint16
	Version uint16
	FileOffset int32
	FileLength int32
}


type StaticPropDictLump struct {
	DictEntries int32
	Name []string // Slice length must equal dictEntries. Validation to be added
}
type StaticPropLeafLump struct {
	LeafEntries int32
	Leaf []uint16 // Slice length must equal leafEntries. Validation to be added
}
type StaticPropLump struct {
	Origin [3]float32
	Angles [3]float32
	PropType uint16
	FirstLeaf uint16
	LeafCount uint16
	Solid string //unsigned char
	Flags string //unsigned char
	Skin int32
	FadeMinDist float32
	FadeMaxDist float32
	LightingOrigin [3]float32
	ForcedFadeScale float32
	MinDXLevel uint16
	MaxDXLevel uint16
	MinCPULevel string //unsigned char
	MaxCPULevel string //unsigned char
	MinGPULevel string //unsigned char
	MaxGPULevel string //unsigned char
	//DiffuseModulation [2]float //dont know if this is correct?
	//Unknown float32
	//DisableXbox360 bool
}
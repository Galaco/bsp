package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

const StaticPropLumpId = 1936749168

type Header struct {
	LumpCount int32
	GameLumps []LumpDef // Slice length must equal lumpCount. Validation to be added
}
func (header Header) SetLumpCount(num int32) Header{
	header.LumpCount = num
	header.GameLumps = make([]LumpDef, header.LumpCount)
	return header
}

type LumpDef struct {
	Id int32
	Flags uint16
	Version uint16
	FileOffset int32
	FileLength int32
}

type GenericGameLump struct {
	Length int32
	Data []byte
}

type StaticPropLump struct {
	DictLump StaticPropDictLump
	LeafLump StaticPropLeafLump
	PropLumps []StaticPropDataLump
}

// Note: Nothing below here is actually implemented, as its primarily game/version specific data.
type StaticPropDictLump struct {
	DictEntries int32
	Name []string // Slice length must equal dictEntries. Validation to be added
}
type StaticPropLeafLump struct {
	LeafEntries int32
	Leaf []uint16 // Slice length must equal leafEntries. Validation to be added
}
type StaticPropDataLump struct {
	Origin mgl32.Vec3
	Angles mgl32.Vec3
	//UniformScale float32 //v11 onwards
	PropType uint16
	FirstLeaf uint16
	LeafCount uint16
	Solid uint8
	Flags uint8
	Skin int32
	FadeMinDist float32
	FadeMaxDist float32
	LightingOrigin mgl32.Vec3
	ForcedFadeScale float32
	MinDXLevel uint16
	MaxDXLevel uint16
	//MinCPULevel uint8
	//MaxCPULevel uint8
	//MinGPULevel uint8
	//MaxGPULevel uint8
	// DiffuseModulation float32 //v7 onwards
	//Unknown float32 //v10 onwards
	//DisableXbox360 bool //v9 onwards
}
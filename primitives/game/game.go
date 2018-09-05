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
	PropLumps []IStaticPropDataLump
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

type IStaticPropDataLump interface {
	GetOrigin() mgl32.Vec3
	GetAngles() mgl32.Vec3
	GetUniformScale() float32 //v11 onwards
	GetPropType() uint16
	GetFirstLeaf() uint16
	GetLeafCount() uint16
	GetSolid() uint8
	GetFlags() uint8
	GetSkin() int32
	GetFadeMinDist() float32
	GetFadeMaxDist() float32
	GetLightingOrigin() mgl32.Vec3
	GetForcedFadeScale() float32
	GetMinDXLevel() uint16
	GetMaxDXLevel() uint16
	GetMinCPULevel() uint8
	GetMaxCPULevel() uint8
	GetMinGPULevel() uint8
	GetMaxGPULevel() uint8
	GetDiffuseModulation() float32 //v7 onwards
	GetUnknown() float32 //v10 onwards
	GetDisableXBox360() bool //v9 onwards
}
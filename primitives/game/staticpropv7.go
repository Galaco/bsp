package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

type StaticPropV7 struct {
	Origin            mgl32.Vec3
	Angles            mgl32.Vec3
	PropType          uint16
	FirstLeaf         uint16
	LeafCount         uint16
	Solid             uint8
	Flags             uint8
	Skin              int32
	FadeMinDist       float32
	FadeMaxDist       float32
	LightingOrigin    mgl32.Vec3
	ForcedFadeScale   float32
	MinDXLevel        uint16
	MaxDXLevel        uint16
	DiffuseModulation float32
}

func (l *StaticPropV7) GetOrigin() mgl32.Vec3 {
	return l.Origin
}

func (l *StaticPropV7) GetAngles() mgl32.Vec3 {
	return l.Angles
}

func (l *StaticPropV7) GetUniformScale() float32 {
	return 1
}

func (l *StaticPropV7) GetPropType() uint16 {
	return l.PropType
}

func (l *StaticPropV7) GetFirstLeaf() uint16 {
	return l.FirstLeaf
}

func (l *StaticPropV7) GetLeafCount() uint16 {
	return l.LeafCount
}

func (l *StaticPropV7) GetSolid() uint8 {
	return l.Solid
}

func (l *StaticPropV7) GetFlags() uint8 {
	return l.Flags
}

func (l *StaticPropV7) GetSkin() int32 {
	return l.Skin
}

func (l *StaticPropV7) GetFadeMinDist() float32 {
	return l.FadeMinDist
}

func (l *StaticPropV7) GetFadeMaxDist() float32 {
	return l.FadeMaxDist
}

func (l *StaticPropV7) GetLightingOrigin() mgl32.Vec3 {
	return l.LightingOrigin
}

func (l *StaticPropV7) GetForcedFadeScale() float32 {
	return l.ForcedFadeScale
}

func (l *StaticPropV7) GetMinDXLevel() uint16 {
	return l.MinDXLevel
}

func (l *StaticPropV7) GetMaxDXLevel() uint16 {
	return l.MaxDXLevel
}

func (l *StaticPropV7) GetMinCPULevel() uint8 {
	return 0
}

func (l *StaticPropV7) GetMaxCPULevel() uint8 {
	return 0
}

func (l *StaticPropV7) GetMinGPULevel() uint8 {
	return 0
}

func (l *StaticPropV7) GetMaxGPULevel() uint8 {
	return 0
}

func (l *StaticPropV7) GetDiffuseModulation() float32 {
	return l.DiffuseModulation
}

func (l *StaticPropV7) GetUnknown() float32 {
	return 0
}

func (l *StaticPropV7) GetDisableXBox360() bool {
	return false
}

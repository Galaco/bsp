package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

type StaticPropV6 struct {
	Origin mgl32.Vec3
	Angles mgl32.Vec3
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
}

func (l *StaticPropV6) GetOrigin() mgl32.Vec3 {
	return l.Origin
}

func (l *StaticPropV6) GetAngles() mgl32.Vec3 {
	return l.Angles
}

func (l *StaticPropV6) GetUniformScale() float32 {
	return 1
}

func (l *StaticPropV6) GetPropType() uint16 {
	return l.PropType
}

func (l *StaticPropV6) GetFirstLeaf() uint16  {
	return l.FirstLeaf
}

func (l *StaticPropV6) GetLeafCount() uint16  {
	return l.LeafCount
}

func (l *StaticPropV6) GetSolid() uint8 {
	return l.Solid
}

func (l *StaticPropV6) GetFlags() uint8 {
	return l.Flags
}

func (l *StaticPropV6) GetSkin() int32 {
	return l.Skin
}

func (l *StaticPropV6) GetFadeMinDist() float32 {
	return l.FadeMinDist
}

func (l *StaticPropV6) GetFadeMaxDist() float32 {
	return l.FadeMaxDist
}

func (l *StaticPropV6) GetLightingOrigin() mgl32.Vec3 {
	return l.LightingOrigin
}

func (l *StaticPropV6) GetForcedFadeScale() float32 {
	return l.ForcedFadeScale
}

func (l *StaticPropV6) GetMinDXLevel() uint16 {
	return l.MinDXLevel
}

func (l *StaticPropV6) GetMaxDXLevel() uint16 {
	return l.MaxDXLevel
}

func (l *StaticPropV6) GetMinCPULevel() uint8 {
	return 0
}

func (l *StaticPropV6) GetMaxCPULevel() uint8 {
	return 0
}

func (l *StaticPropV6) GetMinGPULevel() uint8 {
	return 0
}

func (l *StaticPropV6) GetMaxGPULevel() uint8 {
	return 0
}

func (l *StaticPropV6) GetDiffuseModulation() float32 {
	return 0
}

func (l *StaticPropV6) GetUnknown() float32 {
	return 0
}

func (l *StaticPropV6) GetDisableXBox360() bool {
	return false
}
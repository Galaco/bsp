package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

type StaticPropV8 struct {
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
	MinCPULevel uint8
	MaxCPULevel uint8
	MinGPULevel uint8
	MaxGPULevel uint8
	DiffuseModulation float32
}

func (l *StaticPropV8) GetOrigin() mgl32.Vec3 {
	return l.Origin
}

func (l *StaticPropV8) GetAngles() mgl32.Vec3 {
	return l.Angles
}

func (l *StaticPropV8) GetUniformScale() float32 {
	return 1
}

func (l *StaticPropV8) GetPropType() uint16 {
	return l.PropType
}

func (l *StaticPropV8) GetFirstLeaf() uint16  {
	return l.FirstLeaf
}

func (l *StaticPropV8) GetLeafCount() uint16  {
	return l.LeafCount
}

func (l *StaticPropV8) GetSolid() uint8 {
	return l.Solid
}

func (l *StaticPropV8) GetFlags() uint8 {
	return l.Flags
}

func (l *StaticPropV8) GetSkin() int32 {
	return l.Skin
}

func (l *StaticPropV8) GetFadeMinDist() float32 {
	return l.FadeMinDist
}

func (l *StaticPropV8) GetFadeMaxDist() float32 {
	return l.FadeMaxDist
}

func (l *StaticPropV8) GetLightingOrigin() mgl32.Vec3 {
	return l.LightingOrigin
}

func (l *StaticPropV8) GetForcedFadeScale() float32 {
	return l.ForcedFadeScale
}

func (l *StaticPropV8) GetMinDXLevel() uint16 {
	return 0
}

func (l *StaticPropV8) GetMaxDXLevel() uint16 {
	return 0
}

func (l *StaticPropV8) GetMinCPULevel() uint8 {
	return l.MaxCPULevel
}

func (l *StaticPropV8) GetMaxCPULevel() uint8 {
	return l.MaxCPULevel
}

func (l *StaticPropV8) GetMinGPULevel() uint8 {
	return l.MinGPULevel
}

func (l *StaticPropV8) GetMaxGPULevel() uint8 {
	return l.MaxGPULevel
}

func (l *StaticPropV8) GetDiffuseModulation() float32 {
	return l.DiffuseModulation
}

func (l *StaticPropV8) GetUnknown() float32 {
	return 0
}

func (l *StaticPropV8) GetDisableXBox360() bool {
	return false
}
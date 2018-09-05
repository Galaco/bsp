package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

type StaticPropV9 struct {
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
	DisableXBox360 bool
}

func (l *StaticPropV9) GetOrigin() mgl32.Vec3 {
	return l.Origin
}

func (l *StaticPropV9) GetAngles() mgl32.Vec3 {
	return l.Angles
}

func (l *StaticPropV9) GetUniformScale() float32 {
	return 1
}

func (l *StaticPropV9) GetPropType() uint16 {
	return l.PropType
}

func (l *StaticPropV9) GetFirstLeaf() uint16  {
	return l.FirstLeaf
}

func (l *StaticPropV9) GetLeafCount() uint16  {
	return l.LeafCount
}

func (l *StaticPropV9) GetSolid() uint8 {
	return l.Solid
}

func (l *StaticPropV9) GetFlags() uint8 {
	return l.Flags
}

func (l *StaticPropV9) GetSkin() int32 {
	return l.Skin
}

func (l *StaticPropV9) GetFadeMinDist() float32 {
	return l.FadeMinDist
}

func (l *StaticPropV9) GetFadeMaxDist() float32 {
	return l.FadeMaxDist
}

func (l *StaticPropV9) GetLightingOrigin() mgl32.Vec3 {
	return l.LightingOrigin
}

func (l *StaticPropV9) GetForcedFadeScale() float32 {
	return l.ForcedFadeScale
}

func (l *StaticPropV9) GetMinDXLevel() uint16 {
	return 0
}

func (l *StaticPropV9) GetMaxDXLevel() uint16 {
	return 0
}

func (l *StaticPropV9) GetMinCPULevel() uint8 {
	return l.MaxCPULevel
}

func (l *StaticPropV9) GetMaxCPULevel() uint8 {
	return l.MaxCPULevel
}

func (l *StaticPropV9) GetMinGPULevel() uint8 {
	return l.MinGPULevel
}

func (l *StaticPropV9) GetMaxGPULevel() uint8 {
	return l.MaxGPULevel
}

func (l *StaticPropV9) GetDiffuseModulation() float32 {
	return l.DiffuseModulation
}

func (l *StaticPropV9) GetUnknown() float32 {
	return 0
}

func (l *StaticPropV9) GetDisableXBox360() bool {
	return l.DisableXBox360
}
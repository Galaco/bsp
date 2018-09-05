package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

type StaticPropV11 struct {
	Origin mgl32.Vec3
	Angles mgl32.Vec3
	UniformScale float32
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
	Unknown float32
	DisableXBox360 bool
}

func (l *StaticPropV11) GetOrigin() mgl32.Vec3 {
	return l.Origin
}

func (l *StaticPropV11) GetAngles() mgl32.Vec3 {
	return l.Angles
}

func (l *StaticPropV11) GetUniformScale() float32 {
	return l.UniformScale
}

func (l *StaticPropV11) GetPropType() uint16 {
	return l.PropType
}

func (l *StaticPropV11) GetFirstLeaf() uint16  {
	return l.FirstLeaf
}

func (l *StaticPropV11) GetLeafCount() uint16  {
	return l.LeafCount
}

func (l *StaticPropV11) GetSolid() uint8 {
	return l.Solid
}

func (l *StaticPropV11) GetFlags() uint8 {
	return l.Flags
}

func (l *StaticPropV11) GetSkin() int32 {
	return l.Skin
}

func (l *StaticPropV11) GetFadeMinDist() float32 {
	return l.FadeMinDist
}

func (l *StaticPropV11) GetFadeMaxDist() float32 {
	return l.FadeMaxDist
}

func (l *StaticPropV11) GetLightingOrigin() mgl32.Vec3 {
	return l.LightingOrigin
}

func (l *StaticPropV11) GetForcedFadeScale() float32 {
	return l.ForcedFadeScale
}

func (l *StaticPropV11) GetMinDXLevel() uint16 {
	return 0
}

func (l *StaticPropV11) GetMaxDXLevel() uint16 {
	return 0
}

func (l *StaticPropV11) GetMinCPULevel() uint8 {
	return l.MaxCPULevel
}

func (l *StaticPropV11) GetMaxCPULevel() uint8 {
	return l.MaxCPULevel
}

func (l *StaticPropV11) GetMinGPULevel() uint8 {
	return l.MinGPULevel
}

func (l *StaticPropV11) GetMaxGPULevel() uint8 {
	return l.MaxGPULevel
}

func (l *StaticPropV11) GetDiffuseModulation() float32 {
	return l.DiffuseModulation
}

func (l *StaticPropV11) GetUnknown() float32 {
	return l.Unknown
}

func (l *StaticPropV11) GetDisableXBox360() bool {
	return l.DisableXBox360
}
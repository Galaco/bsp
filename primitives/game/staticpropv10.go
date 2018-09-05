package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

type StaticPropV10 struct {
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
	Unknown float32
	DisableXBox360 bool
	_ [3]byte
}

func (l *StaticPropV10) GetOrigin() mgl32.Vec3 {
	return l.Origin
}

func (l *StaticPropV10) GetAngles() mgl32.Vec3 {
	return l.Angles
}

func (l *StaticPropV10) GetUniformScale() float32 {
	return 1
}

func (l *StaticPropV10) GetPropType() uint16 {
	return l.PropType
}

func (l *StaticPropV10) GetFirstLeaf() uint16  {
	return l.FirstLeaf
}

func (l *StaticPropV10) GetLeafCount() uint16  {
	return l.LeafCount
}

func (l *StaticPropV10) GetSolid() uint8 {
	return l.Solid
}

func (l *StaticPropV10) GetFlags() uint8 {
	return l.Flags
}

func (l *StaticPropV10) GetSkin() int32 {
	return l.Skin
}

func (l *StaticPropV10) GetFadeMinDist() float32 {
	return l.FadeMinDist
}

func (l *StaticPropV10) GetFadeMaxDist() float32 {
	return l.FadeMaxDist
}

func (l *StaticPropV10) GetLightingOrigin() mgl32.Vec3 {
	return l.LightingOrigin
}

func (l *StaticPropV10) GetForcedFadeScale() float32 {
	return l.ForcedFadeScale
}

func (l *StaticPropV10) GetMinDXLevel() uint16 {
	return 0
}

func (l *StaticPropV10) GetMaxDXLevel() uint16 {
	return 0
}

func (l *StaticPropV10) GetMinCPULevel() uint8 {
	return l.MaxCPULevel
}

func (l *StaticPropV10) GetMaxCPULevel() uint8 {
	return l.MaxCPULevel
}

func (l *StaticPropV10) GetMinGPULevel() uint8 {
	return l.MinGPULevel
}

func (l *StaticPropV10) GetMaxGPULevel() uint8 {
	return l.MaxGPULevel
}

func (l *StaticPropV10) GetDiffuseModulation() float32 {
	return l.DiffuseModulation
}

func (l *StaticPropV10) GetUnknown() float32 {
	return l.Unknown
}

func (l *StaticPropV10) GetDisableXBox360() bool {
	return l.DisableXBox360
}
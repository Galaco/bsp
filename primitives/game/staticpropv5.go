package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

type StaticPropV5 struct {
	Origin          mgl32.Vec3
	Angles          mgl32.Vec3
	PropType        uint16
	FirstLeaf       uint16
	LeafCount       uint16
	Solid           uint8
	Flags           uint8
	Skin            int32
	FadeMinDist     float32
	FadeMaxDist     float32
	LightingOrigin  mgl32.Vec3
	ForcedFadeScale float32
}

func (l *StaticPropV5) GetOrigin() mgl32.Vec3 {
	return l.Origin
}

func (l *StaticPropV5) GetAngles() mgl32.Vec3 {
	return l.Angles
}

func (l *StaticPropV5) GetUniformScale() float32 {
	return 1
}

func (l *StaticPropV5) GetPropType() uint16 {
	return l.PropType
}

func (l *StaticPropV5) GetFirstLeaf() uint16 {
	return l.FirstLeaf
}

func (l *StaticPropV5) GetLeafCount() uint16 {
	return l.LeafCount
}

func (l *StaticPropV5) GetSolid() uint8 {
	return l.Solid
}

func (l *StaticPropV5) GetFlags() uint8 {
	return l.Flags
}

func (l *StaticPropV5) GetSkin() int32 {
	return l.Skin
}

func (l *StaticPropV5) GetFadeMinDist() float32 {
	return l.FadeMinDist
}

func (l *StaticPropV5) GetFadeMaxDist() float32 {
	return l.FadeMaxDist
}

func (l *StaticPropV5) GetLightingOrigin() mgl32.Vec3 {
	return l.LightingOrigin
}

func (l *StaticPropV5) GetForcedFadeScale() float32 {
	return l.ForcedFadeScale
}

func (l *StaticPropV5) GetMinDXLevel() uint16 {
	return 0
}

func (l *StaticPropV5) GetMaxDXLevel() uint16 {
	return 0
}

func (l *StaticPropV5) GetMinCPULevel() uint8 {
	return 0
}

func (l *StaticPropV5) GetMaxCPULevel() uint8 {
	return 0
}

func (l *StaticPropV5) GetMinGPULevel() uint8 {
	return 0
}

func (l *StaticPropV5) GetMaxGPULevel() uint8 {
	return 0
}

func (l *StaticPropV5) GetDiffuseModulation() float32 {
	return 0
}

func (l *StaticPropV5) GetUnknown() float32 {
	return 0
}

func (l *StaticPropV5) GetDisableXBox360() bool {
	return false
}

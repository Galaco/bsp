package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

// StaticPropV10 v10 type
// v10 is the model prolific format, most of Valves games use v10
type StaticPropV10 struct {
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
	MinCPULevel       uint8
	MaxCPULevel       uint8
	MinGPULevel       uint8
	MaxGPULevel       uint8
	DiffuseModulation float32
	DisableXBox360    bool
	ExtraFlags        int32
	_ byte
	_ byte
	_ byte
}

// GetOrigin Origin of object in world
func (l *StaticPropV10) GetOrigin() mgl32.Vec3 {
	return l.Origin
}

// GetAngles Rotation of object in world
func (l *StaticPropV10) GetAngles() mgl32.Vec3 {
	return l.Angles
}

// GetUniformScale Uniform scaling of prop (added in this version)
// Not defined in v10
func (l *StaticPropV10) GetUniformScale() float32 {
	return 1
}

// GetPropType prop type
func (l *StaticPropV10) GetPropType() uint16 {
	return l.PropType
}

// GetFirstLeaf Index into StaticPropLeafLump
func (l *StaticPropV10) GetFirstLeaf() uint16 {
	return l.FirstLeaf
}

// GetLeafCount Number of leafs this prop is in
func (l *StaticPropV10) GetLeafCount() uint16 {
	return l.LeafCount
}

// GetSolid is solid
func (l *StaticPropV10) GetSolid() uint8 {
	return l.Solid
}

// GetFlags staticprop flags
func (l *StaticPropV10) GetFlags() uint8 {
	return l.Flags
}

// GetSkin skin index (default 0)
func (l *StaticPropV10) GetSkin() int32 {
	return l.Skin
}

// GetFadeMinDist Distance from prop that it starts to fade
func (l *StaticPropV10) GetFadeMinDist() float32 {
	return l.FadeMinDist
}

// GetFadeMaxDist Distance from prop that it is fully invisible/not rendered
func (l *StaticPropV10) GetFadeMaxDist() float32 {
	return l.FadeMaxDist
}

// GetLightingOrigin world position to sample light from.
// This may differ from prop origin
func (l *StaticPropV10) GetLightingOrigin() mgl32.Vec3 {
	return l.LightingOrigin
}

// GetForcedFadeScale
func (l *StaticPropV10) GetForcedFadeScale() float32 {
	return l.ForcedFadeScale
}

// GetMinDXLevel
// Not defined in v10
func (l *StaticPropV10) GetMinDXLevel() uint16 {
	return 0
}

// GetMaxDXLevel
// Not defined in v10
func (l *StaticPropV10) GetMaxDXLevel() uint16 {
	return 0
}

// GetMinCPULevel minimum cpu to render
func (l *StaticPropV10) GetMinCPULevel() uint8 {
	return l.MinCPULevel
}

// GetMaxCPULevel maximum cpu to render
func (l *StaticPropV10) GetMaxCPULevel() uint8 {
	return l.MaxCPULevel
}

// GetMinGPULevel minimum GPU to render
func (l *StaticPropV10) GetMinGPULevel() uint8 {
	return l.MinGPULevel
}

// GetMaxGPULevel Maximum GPU to render
func (l *StaticPropV10) GetMaxGPULevel() uint8 {
	return l.MaxGPULevel
}

// GetDiffuseModulation
func (l *StaticPropV10) GetDiffuseModulation() float32 {
	return l.DiffuseModulation
}

// GetExtraFlags
func (l *StaticPropV10) GetExtraFlags() int32 {
	return l.ExtraFlags
}

// GetUnknown
func (l *StaticPropV10) GetUnknown() float32 {
	return 0
}

// GetDisableXBox360 should disable on XBox 360
func (l *StaticPropV10) GetDisableXBox360() bool {
	return l.DisableXBox360
}

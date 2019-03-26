package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

// StaticPropV11 v11 type
// v11 introduced uniform staticprop scaling in csgo
// there is trailing [7]byte with unknown purpose right now
type StaticPropV11 struct {
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
	UniformScale      float32
	_                 [7]byte
}

// GetOrigin Origin of object in world
func (l *StaticPropV11) GetOrigin() mgl32.Vec3 {
	return l.Origin
}

// GetAngles Rotation of object in world
func (l *StaticPropV11) GetAngles() mgl32.Vec3 {
	return l.Angles
}

// GetUniformScale Uniform scaling of prop (added in this version)
func (l *StaticPropV11) GetUniformScale() float32 {
	return l.UniformScale
}

// GetPropType prop type
func (l *StaticPropV11) GetPropType() uint16 {
	return l.PropType
}

// GetFirstLeaf Index into StaticPropLeafLump
func (l *StaticPropV11) GetFirstLeaf() uint16 {
	return l.FirstLeaf
}

// GetLeafCount Number of leafs this prop is in
func (l *StaticPropV11) GetLeafCount() uint16 {
	return l.LeafCount
}

// GetSolid is solid
func (l *StaticPropV11) GetSolid() uint8 {
	return l.Solid
}

// GetFlags staticprop flags
func (l *StaticPropV11) GetFlags() uint8 {
	return l.Flags
}

// GetSkin skin index (default 0)
func (l *StaticPropV11) GetSkin() int32 {
	return l.Skin
}

// GetFadeMinDist Distance from prop that it starts to fade
func (l *StaticPropV11) GetFadeMinDist() float32 {
	return l.FadeMinDist
}

// GetFadeMaxDist Distance from prop that it is fully invisible/not rendered
func (l *StaticPropV11) GetFadeMaxDist() float32 {
	return l.FadeMaxDist
}

// GetLightingOrigin World position to sample light from.
// This may differ from prop origin
func (l *StaticPropV11) GetLightingOrigin() mgl32.Vec3 {
	return l.LightingOrigin
}

// GetForcedFadeScale
func (l *StaticPropV11) GetForcedFadeScale() float32 {
	return l.ForcedFadeScale
}

// GetMinDXLevel
// Not defined in v11
func (l *StaticPropV11) GetMinDXLevel() uint16 {
	return 0
}

// GetMaxDXLevel
// Not defined in v11
func (l *StaticPropV11) GetMaxDXLevel() uint16 {
	return 0
}

// GetMinCPULevel minimum cpu to render
func (l *StaticPropV11) GetMinCPULevel() uint8 {
	return l.MinCPULevel
}

// GetMaxCPULevel maximum cpu to render
func (l *StaticPropV11) GetMaxCPULevel() uint8 {
	return l.MaxCPULevel
}

// GetMinGPULevel minimum GPU to render
func (l *StaticPropV11) GetMinGPULevel() uint8 {
	return l.MinGPULevel
}

// GetMaxGPULevel Maximum GPU to render
func (l *StaticPropV11) GetMaxGPULevel() uint8 {
	return l.MaxGPULevel
}

// GetDiffuseModulation
func (l *StaticPropV11) GetDiffuseModulation() float32 {
	return l.DiffuseModulation
}

// GetUnknown
func (l *StaticPropV11) GetUnknown() float32 {
	return 0
}

// GetDisableXBox360 should disable on XBox 360
func (l *StaticPropV11) GetDisableXBox360() bool {
	return l.DisableXBox360
}

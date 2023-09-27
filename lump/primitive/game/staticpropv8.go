package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

// StaticPropV8 v8 type
type StaticPropV8 struct {
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
}

// GetOrigin origin of object in world
func (l *StaticPropV8) GetOrigin() mgl32.Vec3 {
	return l.Origin
}

// GetAngles rotation of object in world
func (l *StaticPropV8) GetAngles() mgl32.Vec3 {
	return l.Angles
}

// GetUniformScale is always 1 (i.e. 1x scale multiplier)
// Not defined in v8
func (l *StaticPropV8) GetUniformScale() float32 {
	return 1
}

// GetPropType prop type
func (l *StaticPropV8) GetPropType() uint16 {
	return l.PropType
}

// GetFirstLeaf index into StaticPropLeafLump
func (l *StaticPropV8) GetFirstLeaf() uint16 {
	return l.FirstLeaf
}

// GetLeafCount number of leafs this prop is in
func (l *StaticPropV8) GetLeafCount() uint16 {
	return l.LeafCount
}

// GetSolid is solid
func (l *StaticPropV8) GetSolid() uint8 {
	return l.Solid
}

// GetFlags prop flags
func (l *StaticPropV8) GetFlags() uint8 {
	return l.Flags
}

// GetSkin skin index of this prop (default 0)
func (l *StaticPropV8) GetSkin() int32 {
	return l.Skin
}

// GetFadeMinDist distance from prop that it starts to fade
func (l *StaticPropV8) GetFadeMinDist() float32 {
	return l.FadeMinDist
}

// GetFadeMaxDist distance from prop that it is fully invisible/not rendered
func (l *StaticPropV8) GetFadeMaxDist() float32 {
	return l.FadeMaxDist
}

// GetLightingOrigin world position to sample light from.
// This may differ from prop origin
func (l *StaticPropV8) GetLightingOrigin() mgl32.Vec3 {
	return l.LightingOrigin
}

// GetForcedFadeScale
func (l *StaticPropV8) GetForcedFadeScale() float32 {
	return l.ForcedFadeScale
}

// GetMinDXLevel Minimum directx level to render this prop
// Not defined in v8
func (l *StaticPropV8) GetMinDXLevel() uint16 {
	return 0
}

// GetMaxDXLevel Maximum directx level to render this prop
// Not defined in v8
func (l *StaticPropV8) GetMaxDXLevel() uint16 {
	return 0
}

// GetMinCPULevel minimum cpu to render
func (l *StaticPropV8) GetMinCPULevel() uint8 {
	return l.MinCPULevel
}

// GetMaxCPULevel maximum cpu to render
func (l *StaticPropV8) GetMaxCPULevel() uint8 {
	return l.MaxCPULevel
}

// GetMinGPULevel minimum GPU to render
func (l *StaticPropV8) GetMinGPULevel() uint8 {
	return l.MinGPULevel
}

// GetMaxGPULevel Maximum GPU to render
func (l *StaticPropV8) GetMaxGPULevel() uint8 {
	return l.MaxGPULevel
}

// GetDiffuseModulation
func (l *StaticPropV8) GetDiffuseModulation() float32 {
	return l.DiffuseModulation
}

// GetUnknown
// Not defined in v8
func (l *StaticPropV8) GetUnknown() float32 {
	return 0
}

// GetDisableXBox360 should disable on XBox 360
// Not defined in v8
func (l *StaticPropV8) GetDisableXBox360() bool {
	return false
}

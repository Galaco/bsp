package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

// StaticPropV4 v4 type
type StaticPropV4 struct {
	Origin         mgl32.Vec3
	Angles         mgl32.Vec3
	PropType       uint16
	FirstLeaf      uint16
	LeafCount      uint16
	Solid          uint8
	Flags          uint8
	Skin           int32
	FadeMinDist    float32
	FadeMaxDist    float32
	LightingOrigin mgl32.Vec3
}

// GetOrigin Origin of object in world
func (l *StaticPropV4) GetOrigin() mgl32.Vec3 {
	return l.Origin
}

// GetAngles rotation of object in world
func (l *StaticPropV4) GetAngles() mgl32.Vec3 {
	return l.Angles
}

// GetUniformScale Uniform scale of object in world
// Not defined in v4
func (l *StaticPropV4) GetUniformScale() float32 {
	return 1
}

// GetPropType prop type
func (l *StaticPropV4) GetPropType() uint16 {
	return l.PropType
}

// GetFirstLeaf Index into StaticPropLeafLump
func (l *StaticPropV4) GetFirstLeaf() uint16 {
	return l.FirstLeaf
}

// GetLeafCount Number of leafs this prop is in
func (l *StaticPropV4) GetLeafCount() uint16 {
	return l.LeafCount
}

// GetSolid is solid
func (l *StaticPropV4) GetSolid() uint8 {
	return l.Solid
}

// GetFlags prop flags
func (l *StaticPropV4) GetFlags() uint8 {
	return l.Flags
}

// GetSkin skin index of this prop (default 0)
func (l *StaticPropV4) GetSkin() int32 {
	return l.Skin
}

// GetFadeMinDist Distance from prop that it starts to fade
func (l *StaticPropV4) GetFadeMinDist() float32 {
	return l.FadeMinDist
}

// GetFadeMaxDist Distance from prop that it is fully invisible/not rendered
func (l *StaticPropV4) GetFadeMaxDist() float32 {
	return l.FadeMaxDist
}

// GetLightingOrigin World position to sample light from.
// This may differ from prop origin
func (l *StaticPropV4) GetLightingOrigin() mgl32.Vec3 {
	return l.LightingOrigin
}

// GetForcedFadeScale
// Not defined in v4
func (l *StaticPropV4) GetForcedFadeScale() float32 {
	return 0
}

// GetMinDXLevel Minimum directx level to render this prop
// Not defined in v4
func (l *StaticPropV4) GetMinDXLevel() uint16 {
	return 0
}

// GetMaxDXLevel Maximum directx level to render this prop
// Not defined in v4
func (l *StaticPropV4) GetMaxDXLevel() uint16 {
	return 0
}

// GetMinCPULevel Minimum CPU type to render this prop
// Not defined in v4
func (l *StaticPropV4) GetMinCPULevel() uint8 {
	return 0
}

// GetMaxCPULevel Maximum CPU type to render this prop
// Not defined in v4
func (l *StaticPropV4) GetMaxCPULevel() uint8 {
	return 0
}

// GetMinGPULevel minimum GPU to render
// Not defined in v4
func (l *StaticPropV4) GetMinGPULevel() uint8 {
	return 0
}

// GetMaxGPULevel Maximum GPU to render
// Not defined in v4
func (l *StaticPropV4) GetMaxGPULevel() uint8 {
	return 0
}

// GetDiffuseModulation
// Not defined in v4
func (l *StaticPropV4) GetDiffuseModulation() float32 {
	return 0
}

// GetUnknown
// Not defined in v4
func (l *StaticPropV4) GetUnknown() float32 {
	return 0
}

// GetDisableXBox360 Should be disabled on xbox 360
// Not defined in v4
func (l *StaticPropV4) GetDisableXBox360() bool {
	return false
}

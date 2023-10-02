package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

// StaticPropV5 v5 type
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

// GetOrigin origin of object in world
func (l *StaticPropV5) GetOrigin() mgl32.Vec3 {
	return l.Origin
}

// GetAngles rotation of object in world
func (l *StaticPropV5) GetAngles() mgl32.Vec3 {
	return l.Angles
}

// GetUniformScale is always 1 (i.e. 1x scale multiplier)
// Not defined in v5
func (l *StaticPropV5) GetUniformScale() float32 {
	return 1
}

// GetPropType prop type
func (l *StaticPropV5) GetPropType() uint16 {
	return l.PropType
}

// GetFirstLeaf index into StaticPropLeafLump
func (l *StaticPropV5) GetFirstLeaf() uint16 {
	return l.FirstLeaf
}

// GetLeafCount number of leafs this prop is in
func (l *StaticPropV5) GetLeafCount() uint16 {
	return l.LeafCount
}

// GetSolid is solid
func (l *StaticPropV5) GetSolid() uint8 {
	return l.Solid
}

// GetFlags prop flags
func (l *StaticPropV5) GetFlags() uint8 {
	return l.Flags
}

// GetSkin skin index of this prop (default 0)
func (l *StaticPropV5) GetSkin() int32 {
	return l.Skin
}

// GetFadeMinDist distance from prop that it starts to fade
func (l *StaticPropV5) GetFadeMinDist() float32 {
	return l.FadeMinDist
}

// GetFadeMaxDist distance from prop that it is fully invisible/not rendered
func (l *StaticPropV5) GetFadeMaxDist() float32 {
	return l.FadeMaxDist
}

// GetLightingOrigin world position to sample light from.
// This may differ from prop origin
func (l *StaticPropV5) GetLightingOrigin() mgl32.Vec3 {
	return l.LightingOrigin
}

// GetForcedFadeScale
func (l *StaticPropV5) GetForcedFadeScale() float32 {
	return l.ForcedFadeScale
}

// GetMinDXLevel Minimum directx level to render this prop
// Not defined in v5
func (l *StaticPropV5) GetMinDXLevel() uint16 {
	return 0
}

// GetMaxDXLevel Maximum directx level to render this prop
// Not defined in v5
func (l *StaticPropV5) GetMaxDXLevel() uint16 {
	return 0
}

// GetMinCPULevel Minimum CPU type to render this prop
// Not defined in v5
func (l *StaticPropV5) GetMinCPULevel() uint8 {
	return 0
}

// GetMaxCPULevel Maximum CPU type to render this prop
// Not defined in v5
func (l *StaticPropV5) GetMaxCPULevel() uint8 {
	return 0
}

// GetMinGPULevel minimum GPU to render
// Not defined in v5
func (l *StaticPropV5) GetMinGPULevel() uint8 {
	return 0
}

// GetMaxGPULevel Maximum GPU to render
// Not defined in v5
func (l *StaticPropV5) GetMaxGPULevel() uint8 {
	return 0
}

// GetDiffuseModulation
// Not defined in v5
func (l *StaticPropV5) GetDiffuseModulation() float32 {
	return 0
}

// GetUnknown
// Not defined in v5
func (l *StaticPropV5) GetUnknown() float32 {
	return 0
}

// GetDisableXBox360 should disable on XBox 360
// Not defined in v5
func (l *StaticPropV5) GetDisableXBox360() bool {
	return false
}

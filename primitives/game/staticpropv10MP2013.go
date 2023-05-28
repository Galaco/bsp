package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

// StaticPropV10MP2013 v6 type
type StaticPropV10MP2013 struct {
	Origin          mgl32.Vec3
	Angles          mgl32.Vec3
	PropType        uint16
	FirstLeaf       uint16
	LeafCount       uint16
	Solid           uint8
	_               byte
	Skin            int32
	FadeMinDist     float32
	FadeMaxDist     float32
	LightingOrigin  mgl32.Vec3
	ForcedFadeScale float32
	MinDXLevel      uint16
	MaxDXLevel      uint16
	Flags           uint8
	_               byte
	LightmapResX    uint16 // lightmap image width
	LightmapResY    uint16 // lightmap image height
	_               [2]byte
}

// GetOrigin origin of object in world
func (l *StaticPropV10MP2013) GetOrigin() mgl32.Vec3 {
	return l.Origin
}

// GetAngles rotation of object in world
func (l *StaticPropV10MP2013) GetAngles() mgl32.Vec3 {
	return l.Angles
}

// GetUniformScale is always 1 (i.e. 1x scale multiplier)
// Not defined in v6
func (l *StaticPropV10MP2013) GetUniformScale() float32 {
	return 1
}

// GetPropType prop type
func (l *StaticPropV10MP2013) GetPropType() uint16 {
	return l.PropType
}

// GetFirstLeaf index into StaticPropLeafLump
func (l *StaticPropV10MP2013) GetFirstLeaf() uint16 {
	return l.FirstLeaf
}

// GetLeafCount number of leafs this prop is in
func (l *StaticPropV10MP2013) GetLeafCount() uint16 {
	return l.LeafCount
}

// GetSolid is solid
func (l *StaticPropV10MP2013) GetSolid() uint8 {
	return l.Solid
}

// GetFlags prop flags
func (l *StaticPropV10MP2013) GetFlags() uint8 {
	return 0
}

// GetSkin skin index of this prop (default 0)
func (l *StaticPropV10MP2013) GetSkin() int32 {
	return l.Skin
}

// GetFadeMinDist distance from prop that it starts to fade
func (l *StaticPropV10MP2013) GetFadeMinDist() float32 {
	return l.FadeMinDist
}

// GetFadeMaxDist distance from prop that it is fully invisible/not rendered
func (l *StaticPropV10MP2013) GetFadeMaxDist() float32 {
	return l.FadeMaxDist
}

// GetLightingOrigin world position to sample light from.
// This may differ from prop origin
func (l *StaticPropV10MP2013) GetLightingOrigin() mgl32.Vec3 {
	return l.LightingOrigin
}

// GetForcedFadeScale
func (l *StaticPropV10MP2013) GetForcedFadeScale() float32 {
	return l.ForcedFadeScale
}

// GetMinDXLevel Minimum directx level to render this prop
func (l *StaticPropV10MP2013) GetMinDXLevel() uint16 {
	return l.MinDXLevel
}

// GetMaxDXLevel Maximum directx level to render this prop
func (l *StaticPropV10MP2013) GetMaxDXLevel() uint16 {
	return l.MaxDXLevel
}

// GetMinCPULevel Minimum CPU type to render this prop
// Not defined in v6
func (l *StaticPropV10MP2013) GetMinCPULevel() uint8 {
	return 0
}

// GetMaxCPULevel maximum cpu to render
// Not defined in v6
func (l *StaticPropV10MP2013) GetMaxCPULevel() uint8 {
	return 0
}

// GetMinGPULevel minimum GPU to render
// Not defined in v6
func (l *StaticPropV10MP2013) GetMinGPULevel() uint8 {
	return 0
}

// GetMaxGPULevel Maximum GPU to render
// Not defined in v6
func (l *StaticPropV10MP2013) GetMaxGPULevel() uint8 {
	return 0
}

// GetDiffuseModulation
// Not defined in v6
func (l *StaticPropV10MP2013) GetDiffuseModulation() float32 {
	return 0
}

// GetUnknown
// Not defined in v6
func (l *StaticPropV10MP2013) GetUnknown() float32 {
	return 0
}

// GetDisableXBox360 should disable on XBox 360
// Not defined in v6
func (l *StaticPropV10MP2013) GetDisableXBox360() bool {
	return false
}

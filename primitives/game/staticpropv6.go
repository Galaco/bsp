package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

// StaticProp v6 type
type StaticPropV6 struct {
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
	MinDXLevel      uint16
	MaxDXLevel      uint16
}

// Origin of object in world
func (l *StaticPropV6) GetOrigin() mgl32.Vec3 {
	return l.Origin
}

// Rotation of object in world
func (l *StaticPropV6) GetAngles() mgl32.Vec3 {
	return l.Angles
}

// Not defined in v6
func (l *StaticPropV6) GetUniformScale() float32 {
	return 1
}

func (l *StaticPropV6) GetPropType() uint16 {
	return l.PropType
}

// Index into StaticPropLeafLump
func (l *StaticPropV6) GetFirstLeaf() uint16 {
	return l.FirstLeaf
}

// Number of leafs this prop is in
func (l *StaticPropV6) GetLeafCount() uint16 {
	return l.LeafCount
}

func (l *StaticPropV6) GetSolid() uint8 {
	return l.Solid
}

func (l *StaticPropV6) GetFlags() uint8 {
	return l.Flags
}

// Skin index of this prop
func (l *StaticPropV6) GetSkin() int32 {
	return l.Skin
}

// Distance from prop that it starts to fade
func (l *StaticPropV6) GetFadeMinDist() float32 {
	return l.FadeMinDist
}

// Distance from prop that it is fully invisible/not rendered
func (l *StaticPropV6) GetFadeMaxDist() float32 {
	return l.FadeMaxDist
}

// World position to sample light from.
// This may differ from prop origin
func (l *StaticPropV6) GetLightingOrigin() mgl32.Vec3 {
	return l.LightingOrigin
}

func (l *StaticPropV6) GetForcedFadeScale() float32 {
	return l.ForcedFadeScale
}

// Minimum directx level to render this prop
func (l *StaticPropV6) GetMinDXLevel() uint16 {
	return l.MinDXLevel
}

// Maximum directx level to render this prop
func (l *StaticPropV6) GetMaxDXLevel() uint16 {
	return l.MaxDXLevel
}

// Minimum CPU type to render this prop
// Not defined in v6
func (l *StaticPropV6) GetMinCPULevel() uint8 {
	return 0
}

// Maximum CPU type to render this prop
// Not defined in v6
func (l *StaticPropV6) GetMaxCPULevel() uint8 {
	return 0
}

// Not defined in v6
func (l *StaticPropV6) GetMinGPULevel() uint8 {
	return 0
}

// Not defined in v6
func (l *StaticPropV6) GetMaxGPULevel() uint8 {
	return 0
}

// Not defined in v6
func (l *StaticPropV6) GetDiffuseModulation() float32 {
	return 0
}

// Not defined in v6
func (l *StaticPropV6) GetUnknown() float32 {
	return 0
}

// Not defined in v6
func (l *StaticPropV6) GetDisableXBox360() bool {
	return false
}

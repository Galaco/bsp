package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

// StaticProp v5 type
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

// Origin of object in world
func (l *StaticPropV5) GetOrigin() mgl32.Vec3 {
	return l.Origin
}

// Rotation of object in world
func (l *StaticPropV5) GetAngles() mgl32.Vec3 {
	return l.Angles
}

// Not defined in v5
func (l *StaticPropV5) GetUniformScale() float32 {
	return 1
}

func (l *StaticPropV5) GetPropType() uint16 {
	return l.PropType
}

// Index into StaticPropLeafLump
func (l *StaticPropV5) GetFirstLeaf() uint16 {
	return l.FirstLeaf
}

// Number of leafs this prop is in
func (l *StaticPropV5) GetLeafCount() uint16 {
	return l.LeafCount
}

func (l *StaticPropV5) GetSolid() uint8 {
	return l.Solid
}

func (l *StaticPropV5) GetFlags() uint8 {
	return l.Flags
}

// Skin index of this prop
func (l *StaticPropV5) GetSkin() int32 {
	return l.Skin
}

// Distance from prop that it starts to fade
func (l *StaticPropV5) GetFadeMinDist() float32 {
	return l.FadeMinDist
}

// Distance from prop that it is fully invisible/not rendered
func (l *StaticPropV5) GetFadeMaxDist() float32 {
	return l.FadeMaxDist
}

// World position to sample light from.
// This may differ from prop origin
func (l *StaticPropV5) GetLightingOrigin() mgl32.Vec3 {
	return l.LightingOrigin
}

func (l *StaticPropV5) GetForcedFadeScale() float32 {
	return l.ForcedFadeScale
}

// Minimum directx level to render this prop
// Not defined in v5
func (l *StaticPropV5) GetMinDXLevel() uint16 {
	return 0
}

// Maximum directx level to render this prop
// Not defined in v5
func (l *StaticPropV5) GetMaxDXLevel() uint16 {
	return 0
}

// Minimum CPU type to render this prop
// Not defined in v5
func (l *StaticPropV5) GetMinCPULevel() uint8 {
	return 0
}

// Maximum CPU type to render this prop
// Not defined in v5
func (l *StaticPropV5) GetMaxCPULevel() uint8 {
	return 0
}

// Not defined in v5
func (l *StaticPropV5) GetMinGPULevel() uint8 {
	return 0
}

// Not defined in v5
func (l *StaticPropV5) GetMaxGPULevel() uint8 {
	return 0
}

// Not defined in v5
func (l *StaticPropV5) GetDiffuseModulation() float32 {
	return 0
}

// Not defined in v5
func (l *StaticPropV5) GetUnknown() float32 {
	return 0
}

// Not defined in v5
func (l *StaticPropV5) GetDisableXBox360() bool {
	return false
}

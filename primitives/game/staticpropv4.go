package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

// StaticProp v4 type
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

// Origin of object in world
func (l *StaticPropV4) GetOrigin() mgl32.Vec3 {
	return l.Origin
}

// Rotation of object in world
func (l *StaticPropV4) GetAngles() mgl32.Vec3 {
	return l.Angles
}

// Uniform scale of object in world
// Not defined in v4
func (l *StaticPropV4) GetUniformScale() float32 {
	return 1
}

func (l *StaticPropV4) GetPropType() uint16 {
	return l.PropType
}

// Index into StaticPropLeafLump
func (l *StaticPropV4) GetFirstLeaf() uint16 {
	return l.FirstLeaf
}

// Number of leafs this prop is in
func (l *StaticPropV4) GetLeafCount() uint16 {
	return l.LeafCount
}

func (l *StaticPropV4) GetSolid() uint8 {
	return l.Solid
}

func (l *StaticPropV4) GetFlags() uint8 {
	return l.Flags
}

// Skin index of this prop
func (l *StaticPropV4) GetSkin() int32 {
	return l.Skin
}

// Distance from prop that it starts to fade
func (l *StaticPropV4) GetFadeMinDist() float32 {
	return l.FadeMinDist
}

// Distance from prop that it is fully invisible/not rendered
func (l *StaticPropV4) GetFadeMaxDist() float32 {
	return l.FadeMaxDist
}

// World position to sample light from.
// This may differ from prop origin
func (l *StaticPropV4) GetLightingOrigin() mgl32.Vec3 {
	return l.LightingOrigin
}

// Not defined in v4
func (l *StaticPropV4) GetForcedFadeScale() float32 {
	return 0
}

// Minimum directx level to render this prop
// Not defined in v4
func (l *StaticPropV4) GetMinDXLevel() uint16 {
	return 0
}

// Maximum directx level to render this prop
// Not defined in v4
func (l *StaticPropV4) GetMaxDXLevel() uint16 {
	return 0
}

// Minimum CPU type to render this prop
// Not defined in v4
func (l *StaticPropV4) GetMinCPULevel() uint8 {
	return 0
}

// Maximum CPU type to render this prop
// Not defined in v4
func (l *StaticPropV4) GetMaxCPULevel() uint8 {
	return 0
}

// Not defined in v4
func (l *StaticPropV4) GetMinGPULevel() uint8 {
	return 0
}

// Not defined in v4
func (l *StaticPropV4) GetMaxGPULevel() uint8 {
	return 0
}

// Not defined in v4
func (l *StaticPropV4) GetDiffuseModulation() float32 {
	return 0
}

// Not defined in v4
func (l *StaticPropV4) GetUnknown() float32 {
	return 0
}

// Should be disabled on xbox 360
// Not defined in v4
func (l *StaticPropV4) GetDisableXBox360() bool {
	return false
}

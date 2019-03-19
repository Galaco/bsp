package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

// StaticProp v7 type
type StaticPropV7 struct {
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
	MinDXLevel        uint16
	MaxDXLevel        uint16
	DiffuseModulation float32
}

// Origin of object in world
func (l *StaticPropV7) GetOrigin() mgl32.Vec3 {
	return l.Origin
}

// Rotation of object in world
func (l *StaticPropV7) GetAngles() mgl32.Vec3 {
	return l.Angles
}

// Not defined in v7
func (l *StaticPropV7) GetUniformScale() float32 {
	return 1
}

func (l *StaticPropV7) GetPropType() uint16 {
	return l.PropType
}

// Index into StaticPropLeafLump
func (l *StaticPropV7) GetFirstLeaf() uint16 {
	return l.FirstLeaf
}

// Number of leafs this prop is in
func (l *StaticPropV7) GetLeafCount() uint16 {
	return l.LeafCount
}

func (l *StaticPropV7) GetSolid() uint8 {
	return l.Solid
}

func (l *StaticPropV7) GetFlags() uint8 {
	return l.Flags
}

// Skin index of this prop
func (l *StaticPropV7) GetSkin() int32 {
	return l.Skin
}

// Distance from prop that it starts to fade
func (l *StaticPropV7) GetFadeMinDist() float32 {
	return l.FadeMinDist
}

// Distance from prop that it is fully invisible/not rendered
func (l *StaticPropV7) GetFadeMaxDist() float32 {
	return l.FadeMaxDist
}

// World position to sample light from.
// This may differ from prop origin
func (l *StaticPropV7) GetLightingOrigin() mgl32.Vec3 {
	return l.LightingOrigin
}

func (l *StaticPropV7) GetForcedFadeScale() float32 {
	return l.ForcedFadeScale
}

func (l *StaticPropV7) GetMinDXLevel() uint16 {
	return l.MinDXLevel
}

func (l *StaticPropV7) GetMaxDXLevel() uint16 {
	return l.MaxDXLevel
}

// Not defined in v7
func (l *StaticPropV7) GetMinCPULevel() uint8 {
	return 0
}

// Not defined in v7
func (l *StaticPropV7) GetMaxCPULevel() uint8 {
	return 0
}

// Not defined in v7
func (l *StaticPropV7) GetMinGPULevel() uint8 {
	return 0
}

// Not defined in v7
func (l *StaticPropV7) GetMaxGPULevel() uint8 {
	return 0
}

func (l *StaticPropV7) GetDiffuseModulation() float32 {
	return l.DiffuseModulation
}

// Not defined in v7
func (l *StaticPropV7) GetUnknown() float32 {
	return 0
}

// Not defined in v7
func (l *StaticPropV7) GetDisableXBox360() bool {
	return false
}

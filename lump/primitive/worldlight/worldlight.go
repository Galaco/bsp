package worldlight

import (
	"github.com/go-gl/mathgl/mgl32"
)

// EmitType light emission mode.
type EmitType uint8 // assumed this is 1 byte..

const (
	// EmitSurface 90 degree spotlight
	EmitSurface EmitType = 0
	// EmitPoint simple point light source
	EmitPoint EmitType = 1
	// EmitSpotLight spotlight with penumbra
	EmitSpotLight EmitType = 2
	// EmitSkyLight directional light with no falloff (surface must trace to SKY texture)
	EmitSkyLight EmitType = 3
	// EmitQuakeLight linear falloff, non-lambertian
	EmitQuakeLight EmitType = 4
	// EmitSkyAmbient spherical light source with no falloff (surface must trace to SKY texture)
	EmitSkyAmbient EmitType = 5
)

// WorldLight is a single light in the world
// This data may also be stored in entdata
type WorldLight struct {
	// Origin is position in the world
	Origin mgl32.Vec3 `json:"origin"`
	// Intensity
	Intensity mgl32.Vec3 `json:"intensity"`
	// Normal
	Normal mgl32.Vec3 `json:"normal"`
	// Cluster
	Cluster int32 `json:"cluster"`
	// Type
	Type EmitType `json:"type"`
	// Unknown1
	// @TODO Think for alignments sake with is uint8. May be 3 bytes padding...
	Unknown1 [3]byte `json:"unknown1"`
	// Style
	Style int32 `json:"style"`
	// StopDot
	StopDot float32 `json:"stopDot"`
	// StopDot2
	StopDot2 float32 `json:"stopDot2"`
	// Exponent
	Exponent float32 `json:"exponent"`
	// Radius
	Radius float32 `json:"radius"`
	// ConstantAttenuation
	ConstantAttenuation float32 `json:"constantAttenuation"`
	// LinearAttenuation
	LinearAttenuation float32 `json:"linearAttenuation"`
	// QuadraticAttenuation
	QuadraticAttenuation float32 `json:"quadraticAttenuation"`
	// Flags
	Flags int32 `json:"flags"`
	// TexInfo
	TexInfo int32 `json:"texInfo"`
	// Owner
	Owner int32 `json:"owner"`
}

// WorldLightHDR is a single light in the world
// This data may also be stored in entdata
type WorldLightHDR struct {
	// Origin is position in the world
	Origin mgl32.Vec3 `json:"origin" bsp:"all"`
	// Intensity
	Intensity mgl32.Vec3 `json:"intensity" bsp:"all"`
	// Normal
	Normal mgl32.Vec3 `json:"normal" bsp:"all"`
	// Cluster
	Cluster int32 `json:"cluster" bsp:"all"`
	// Type
	Type EmitType `json:"type" bsp:"all"`
	// Unknown1
	// @TODO Think for alignments sake with is uint8. May be 3 bytes padding...
	Unknown1 [3]byte `json:"unknown1" bsp:"all"`
	// Style
	Style int32 `json:"style" bsp:"all"`
	// StopDot
	StopDot float32 `json:"stopDot" bsp:"all"`
	// StopDot2
	StopDot2 float32 `json:"stopDot2" bsp:"all"`
	// Exponent
	Exponent float32 `json:"exponent" bsp:"all"`
	// Radius
	Radius float32 `json:"radius" bsp:"all"`
	// ConstantAttenuation
	ConstantAttenuation float32 `json:"constantAttenuation" bsp:"all"`
	// LinearAttenuation
	LinearAttenuation float32 `json:"linearAttenuation" bsp:"all"`
	// QuadraticAttenuation
	QuadraticAttenuation float32 `json:"quadraticAttenuation" bsp:"all"`
	// Flags
	Flags int32 `json:"flags" bsp:"all"`
	// TexInfo
	TexInfo int32 `json:"texInfo" bsp:"all"`
	// Owner
	Owner int32 `json:"owner" bsp:"all"`
	// @TODO We don't know what this is.
	Unknown2 [12]byte `json:"unknown2" bsp:"v21"`
	// @TODO We don't know what this is, but v20 has it and v21 does not.
	//Unknown3 [8]byte `json:"unknown3" bsp:"v20"`
}

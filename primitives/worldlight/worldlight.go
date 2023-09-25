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
	Origin mgl32.Vec3
	// Intensity
	Intensity mgl32.Vec3
	// Normal
	Normal mgl32.Vec3
	// Cluster
	Cluster int32
	// Type
	Type EmitType
	// Unknown1
	// @TODO Think for alignments sake with is uint8. May be 3 bytes padding...
	Unknown1 [3]byte
	// Style
	Style int32
	// Stopdot
	Stopdot float32
	// Stopdot2
	Stopdot2 float32
	// Exponent
	Exponent float32
	// Radius
	Radius float32
	// ConstantAttenuation
	ConstantAttenuation float32
	// LinearAttenuation
	LinearAttenuation float32
	// QuadraticAttenuation
	QuadraticAttenuation float32
	// Flags
	Flags int32
	// TexInfo
	TexInfo int32
	// Owner
	Owner int32
}

// WorldLightHDR is a single light in the world
// This data may also be stored in entdata
type WorldLightHDR struct {
	// Origin is position in the world
	Origin mgl32.Vec3
	// Intensity
	Intensity mgl32.Vec3
	// Normal
	Normal mgl32.Vec3
	// Cluster
	Cluster int32
	// Type
	Type EmitType
	// Unknown1
	// @TODO Think for alignments sake with is uint8. May be 3 bytes padding...
	Unknown1 [3]byte
	// Style
	Style int32
	// Stopdot
	Stopdot float32
	// Stopdot2
	Stopdot2 float32
	// Exponent
	Exponent float32
	// Radius
	Radius float32
	// ConstantAttenuation
	ConstantAttenuation float32
	// LinearAttenuation
	LinearAttenuation float32
	// QuadraticAttenuation
	QuadraticAttenuation float32
	// Flags
	Flags int32
	// TexInfo
	TexInfo int32
	// Owner
	Owner int32
	// Unknown2
	Unknown2 [12]byte
}

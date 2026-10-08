package lumps

import "github.com/go-gl/mathgl/mgl32"

// EmitType is the signed 32-bit on-disk light emission mode. Unknown values
// remain representable and are preserved without reinterpretation.
type EmitType int32

const (
	// EmitSurface is a 90 degree spotlight.
	EmitSurface EmitType = iota
	// EmitPoint is a simple point light source.
	EmitPoint
	// EmitSpotLight is a spotlight with penumbra.
	EmitSpotLight
	// EmitSkyLight is directional light with no falloff.
	EmitSkyLight
	// EmitQuakeLight has linear falloff and is non-lambertian.
	EmitQuakeLight
	// EmitSkyAmbient is spherical light with no falloff.
	EmitSkyAmbient
)

// WorldLight is a worldlight record in either the LDR or HDR channel.
// ShadowCastOffset is present only in version 1; version 0 requires zero bits.
type WorldLight struct {
	Origin               mgl32.Vec3
	Intensity            mgl32.Vec3
	Normal               mgl32.Vec3
	ShadowCastOffset     mgl32.Vec3
	Cluster              int32
	Type                 EmitType
	Style                int32
	Stopdot              float32
	Stopdot2             float32
	Exponent             float32
	Radius               float32
	ConstantAttenuation  float32
	LinearAttenuation    float32
	QuadraticAttenuation float32
	Flags                int32
	TexInfo              int32
	Owner                int32
}

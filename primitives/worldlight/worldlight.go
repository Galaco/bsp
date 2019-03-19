package worldlight

import (
	"github.com/go-gl/mathgl/mgl32"
)

type EmitType uint8 // assumed this is 1 byte..

const EMIT_SURFACE EmitType = 0    // 90 degree spotlight
const EMIT_POINT EmitType = 1      // simple point light source
const EMIT_SPOTLIGHT EmitType = 2  // spotlight with penumbra
const EMIT_SKYLIGHT EmitType = 3   // directional light with no falloff (surface must trace to SKY texture)
const EMIT_QUAKELIGHT EmitType = 4 // linear falloff, non-lambertian
const EMIT_SKYAMBIENT EmitType = 5 // spherical light source with no falloff (surface must trace to SKY texture)

// A single light in the world
// This data may also be stored in entdata
type WorldLight struct {
	Origin               mgl32.Vec3
	Intensity            mgl32.Vec3
	Normal               mgl32.Vec3
	Cluster              int32
	Type                 EmitType //Think for alignments sake with is uint8. May be 3 bytes padding...
	_                    [3]byte
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

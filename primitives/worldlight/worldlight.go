package worldlight

import (
	"github.com/go-gl/mathgl/mgl32"
)

/* Assuming this is 8bits
enum emittype_t
{
	emit_surface,		// 90 degree spotlight
	emit_point,			// simple point light source
	emit_spotlight,		// spotlight with penumbra
	emit_skylight,		// directional light with no falloff (surface must trace to SKY texture)
	emit_quakelight,	// linear falloff, non-lambertian
	emit_skyambient,	// spherical light source with no falloff (surface must trace to SKY texture)
};
 */

type EmitType uint8

const EMIT_SURFACE EmitType = 0
const EMIT_POINT EmitType = 1
const EMIT_POINTLIGHT EmitType = 2
const EMIT_SKYLIGHT EmitType = 3
const EMIT_QUAKELIGHT EmitType = 4
const EMIT_SKYAMBIENT EmitType = 5

type WorldLight struct {
	Origin mgl32.Vec3
	Intensity mgl32.Vec3
	Normal mgl32.Vec3
	Cluster int32
	Type EmitType //Think for alignments sake with is uint8. May be 3 bytes padding...
	_ [3]byte
	Style int32
	Stopdot float32
	Stopdot2 float32
	Exponent float32
	Radius float32
	ConstantAttenuation float32
	LinearAttenuation float32
	QuadraticAttenuation float32
	Flags int32
	TexInfo int32
	Owner int32
}
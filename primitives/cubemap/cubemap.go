package cubemap

import (
	"github.com/go-gl/mathgl/mgl32"
)

type CubemapSample struct {
	Origin mgl32.Vec3
	Size byte
	AlignmentPadding [3]byte
}
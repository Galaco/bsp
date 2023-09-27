package cubemap

import (
	"github.com/go-gl/mathgl/mgl32"
)

// CubemapSample
type CubemapSample struct {
	// Origin is sample location
	Origin mgl32.Vec3 `json:"origin"`
	// Size
	Size byte `json:"size"`
	// AlignmentPadding - probably unused, likely exists for ensure 4byte alignment during read/write.
	AlignmentPadding [3]byte `json:"alignmentPadding"`
}

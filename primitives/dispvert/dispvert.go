package dispvert

import (
	"github.com/go-gl/mathgl/mgl32"
)

// DispVert represents a single vertex on a displacement face. It defines properties
// required to compute the correct position of the displacement vertex after transformation.
type DispVert struct {
	// Vec is direction of vertex from its default computed position
	Vec mgl32.Vec3
	// Dist is distance of vertex from its default computed position
	Dist float32
	// Alpha is the alpha value of this vertex (normally used for blend materials)
	Alpha float32
}

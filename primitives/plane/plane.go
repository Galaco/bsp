package plane

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Plane struct {
	Normal   mgl32.Vec3 // normal vector
	Distance float32    // distance from origin
	AxisType int32      // plane axis identifier
}

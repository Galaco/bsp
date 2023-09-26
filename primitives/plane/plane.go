package plane

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Plane represents an arbitrary plane of infinite length.
type Plane struct {
	// Normal is normal vector
	Normal mgl32.Vec3 `json:"normal"` // normal vector
	// Distance is distance from origin
	Distance float32 `json:"distance"`
	// AxisType is plane axis identifier
	AxisType int32 `json:"axisType"`
}

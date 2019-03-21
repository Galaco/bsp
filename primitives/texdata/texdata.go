package texdata

import (
	"github.com/go-gl/mathgl/mgl32"
)

// TexData contains texture properties
type TexData struct {
	// Reflectivity is reflectivity vector
	Reflectivity mgl32.Vec3
	// NameStringTableID
	NameStringTableID int32
	// Width is texture width
	Width int32
	// Height is texture height
	Height int32
	// ViewWidth
	ViewWidth int32
	// ViewHeight
	ViewHeight int32
}

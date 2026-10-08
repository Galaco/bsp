package lumps

import (
	"github.com/go-gl/mathgl/mgl32"
)

// TextureData contains texture properties
type TextureData struct {
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

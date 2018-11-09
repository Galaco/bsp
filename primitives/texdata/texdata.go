package texdata

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Texture properties
type TexData struct {
	Reflectivity      mgl32.Vec3
	NameStringTableID int32
	Width             int32
	Height            int32
	ViewWidth         int32
	ViewHeight        int32
}

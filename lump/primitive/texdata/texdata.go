package texdata

import (
	"github.com/go-gl/mathgl/mgl32"
)

// TexData contains texture properties
type TexData struct {
	// Reflectivity is reflectivity vector
	Reflectivity mgl32.Vec3 `json:"reflectivity"`
	// NameStringTableID
	NameStringTableID int32 `json:"nameStringTableId"`
	// Width is texture width
	Width int32 `json:"width"`
	// Height is texture height
	Height int32 `json:"height"`
	// ViewWidth
	ViewWidth int32 `json:"viewWidth"`
	// ViewHeight
	ViewHeight int32 `json:"viewHeight"`
}

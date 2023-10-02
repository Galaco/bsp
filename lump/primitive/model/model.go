package model

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Model is a BSP root node.
type Model struct {
	// Mins is bounding volume smallest value corner
	Mins mgl32.Vec3 `json:"mins"`
	// Maxs is bounding volume largest value corner
	Maxs mgl32.Vec3 `json:"maxs"`
	// Origin is the center/specified origin position
	Origin mgl32.Vec3 `json:"origin"`
	// HeadNode
	HeadNode int32 `json:"headNode"`
	// FirstFace is index into the faces lump
	FirstFace int32 `json:"firstFace"`
	// NumFaces is number of faces in this model
	NumFaces int32 `json:"numFaces"`
}

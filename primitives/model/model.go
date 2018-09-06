package model

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Model struct {
	Mins      mgl32.Vec3
	Maxs      mgl32.Vec3
	Origin    mgl32.Vec3
	HeadNode  int32
	FirstFace int32
	NumFaces  int32
}

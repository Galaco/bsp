package dispinfo

import (
	"github.com/go-gl/mathgl/mgl32"
)

const MaxDispCornerNeightbours = 4

type DispInfo struct {
	StartPosition mgl32.Vec3
	DispVertStart int32
	DispTriStart int32
	Power int32
	MinTess int32
	SmoothingAngle float32
	Contents int32
	MapFace uint16
	_ [2]byte
	LightmapAlphaStart int32
	LightmapSampleStartPosition int32
	Ignore [32]uint32
	//EdgeNeighbors [4]DispNeighbor
	//CornerNeighbors [4]DispCornerNeighbors
	//AllowedVerts [8]uint32

}

type DispNeighbor struct {
	SubNeighbours [2]DispSubNeighbor
}

type DispSubNeighbor struct {
	Index uint16 // 0xFFFF if no neighbor
	NeighborOrientation uint8
	Span uint8
	NeighborSpan uint8
}

type DispCornerNeighbors struct {
	Neighbors [MaxDispCornerNeightbours]uint16
	NumNeighbors uint8
}
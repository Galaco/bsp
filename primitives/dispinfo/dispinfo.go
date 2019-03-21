package dispinfo

import (
	"github.com/go-gl/mathgl/mgl32"
)

// MaxDispCornerNeighbours maximum number of adjoining neigherbours per corner.
const MaxDispCornerNeighbours = 4

// DispInfo
type DispInfo struct {
	// StartPosition
	StartPosition mgl32.Vec3
	// DispVertStart
	DispVertStart int32
	// DispTriStart
	DispTriStart int32
	// Power
	Power int32
	// MinTess
	MinTess int32
	// SmoothingAngle
	SmoothingAngle float32
	// Contents
	Contents int32
	// MapFace
	MapFace uint16
	_       [2]byte
	// LightmapAlphaStart
	LightmapAlphaStart int32
	// LightmapSampleStartPosition
	LightmapSampleStartPosition int32
	// Ignore contains bytes with unknown purpose and representation. This should be updated as purpose id discovered
	Ignore [32]uint32
	//EdgeNeighbors [4]DispNeighbor
	//CornerNeighbors [4]DispCornerNeighbors
	//AllowedVerts [8]uint32

}

// DispNeighbor
type DispNeighbor struct {
	// SubNeighbours
	SubNeighbours [2]DispSubNeighbor
}

// DispSubNeighbor
type DispSubNeighbor struct {
	// Index
	// 0xFFFF if no neighbor
	Index uint16
	// NeighborOrientation
	NeighborOrientation uint8
	// Span
	Span uint8
	// NeighborSpan
	NeighborSpan uint8
}

// DispCornerNeighbors
type DispCornerNeighbors struct {
	// Neighbors
	Neighbors [MaxDispCornerNeighbours]uint16
	// NumNeighbors is number of neighbours
	NumNeighbors uint8
}

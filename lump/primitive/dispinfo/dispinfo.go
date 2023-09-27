package dispinfo

import (
	"github.com/go-gl/mathgl/mgl32"
)

// MaxDispCornerNeighbours maximum number of adjoining neigherbours per corner.
const MaxDispCornerNeighbours = 4

// DispInfo
type DispInfo struct {
	// StartPosition
	StartPosition mgl32.Vec3 `json:"startPosition"`
	// DispVertStart
	DispVertStart int32 `json:"dispVertStart"`
	// DispTriStart
	DispTriStart int32 `json:"dispTriStart"`
	// Power
	Power int32 `json:"power"`
	// MinTess
	MinTess int32 `json:"minTess"`
	// SmoothingAngle
	SmoothingAngle float32 `json:"smoothingAngle"`
	// Contents
	Contents int32 `json:"contents"`
	// MapFace
	MapFace uint16 `json:"mapFace"`
	// Unknown1 are unknown bytes.
	// @TODO figure these out.
	Unknown1 [2]byte `json:"unknown1"`
	// LightmapAlphaStart
	LightmapAlphaStart int32 `json:"lightmapAlphaStart"`
	// LightmapSampleStartPosition
	LightmapSampleStartPosition int32 `json:"lightmapSampleStartPosition"`
	// Unknown2 contains bytes with unknown purpose and representation.
	// @TODO This should be updated as purpose id discovered
	Unknown2 [32]uint32 `json:"unknown2"`
	//EdgeNeighbors [4]DispNeighbor
	//CornerNeighbors [4]DispCornerNeighbors
	//AllowedVerts [8]uint32

}

// DispNeighbor
type DispNeighbor struct {
	// SubNeighbours
	SubNeighbours [2]DispSubNeighbor `json:"subNeighbours"`
}

// DispSubNeighbor
type DispSubNeighbor struct {
	// Index
	// 0xFFFF if no neighbor
	Index uint16 `json:"index"`
	// NeighborOrientation
	NeighborOrientation uint8 `json:"neighborOrientation"`
	// Span
	Span uint8 `json:"span"`
	// NeighborSpan
	NeighborSpan uint8 `json:"neighborSpan"`
}

// DispCornerNeighbors
type DispCornerNeighbors struct {
	// Neighbors
	Neighbors [MaxDispCornerNeighbours]uint16 `json:"neighbors"`
	// NumNeighbors is number of neighbours
	NumNeighbors uint8 `json:"numNeighbors"`
}

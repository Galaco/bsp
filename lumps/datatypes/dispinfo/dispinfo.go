package dispinfo

import "github.com/galaco/bsp/lumps/datatypes/common"

type DispInfo struct {
	StartPosition common.Vector
	DispVertStart int32
	DispTriStart int32
	Power int32
	MinTess int32
	SmoothingAngle float32
	Contents int32
	MapFace uint16
	_ uint16
	LightmapAlphaStart int32
	LightmapSampleStartPosition int32
	EdgeNeighbors [4]DispNeighbor
	CornerNeighbors [4]DispCornerNeighbors
	_ int32
	AllowedVerts [10]uint32
}

type DispNeighbor struct {
	SubNeighbours [2]DispSubNeighbor
}

type DispSubNeighbor struct {
	Index uint16 // 0xFFFF if no neighbor
	NeighborOrientation byte
	Span byte
	NeighborSpan byte
	//_ [3]byte
}

type DispCornerNeighbors struct {
	Neighbors [4]uint16
	NumNeighbors byte
	//_ [3]byte
}
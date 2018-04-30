package common

import "github.com/galaco/bsp/primitives/face"

type Vector struct {
	X float32
	Y float32
	Z float32
}

type VectorInt32 struct {
	X int32
	Y int32
	Z int32
}

type ColorRGBExponent32 struct {
	R byte
	G byte
	B byte
	Exponent byte
}

const MAX_POINTS_ON_FIXED_WINDING = 12
//const MAX_DISPVERTS = 0
//const MAX_DISPTRIS = 0


type Winding struct {
	Original int32 // qboolean = int32
	NumPoints int32
	Points []Vector
}

type Side struct {
	PlaneNum int32
	TexInfo int32
	MapDisp *MapDispInfo
	Winding *Winding
	Original *Side
	Contents int32
	Surf int32
	Visible int32
	Tested int32
	Bevel int32
	Next *Side
	OrigIndex int32
	Id int32
	SmoothingGroups uint32
	AOverlayIds []int32
	AWaterOverlayIds []int32
	DynamicShadowsEnabled bool
}

type MapDispInfo struct {
	Face face.Face
	EntityNum int32
	Power int32
	MinTess int32
	SmoothingAngle float32
	uAxis Vector
	vAxis Vector
	StartPosition Vector
	AlphaValues []float32
	MaxDispDist float32
	DispDists []float32
	VectorDisps []Vector
	VectorOffsets []Vector
	Contents int32
	BrushSideID int32
	Flags int32

	// #ifdef VSVMFIO
	// Elevation float32
	// OffsetNormals []Vector
}
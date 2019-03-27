package common

import (
	"github.com/galaco/bsp/primitives/face"
	"github.com/go-gl/mathgl/mgl32"
)

// ColorRGBExponent32 is a 4 byte color representation (RGBExp)
type ColorRGBExponent32 struct {
	// R red
	R uint8
	// G green
	G uint8
	// B blue
	B uint8
	// Exponent exponent
	Exponent int8
}

// Winding
type Winding struct {
	// Original
	Original int32 // qboolean = int32
	// NumPoints
	NumPoints int32
	// Points
	Points []mgl32.Vec3
}

// Side
type Side struct {
	// PlaneNum
	PlaneNum int32
	// TexInfo
	TexInfo int32
	// MapDisp
	MapDisp *MapDispInfo
	// Winding
	Winding *Winding
	// Original
	Original *Side
	// Contents
	Contents int32
	// Surf
	Surf int32
	// Visible
	Visible int32
	// Tested
	Tested int32
	// Bevel
	Bevel int32
	// Next
	Next *Side
	// OrigIndex
	OrigIndex int32
	// Id is side Id
	Id int32
	// SmoothingGroups
	SmoothingGroups uint32
	// AOverlayIds
	AOverlayIds []int32
	// AWaterOverlayIds
	AWaterOverlayIds []int32
	// DynamicShadowsEnabled are dynamic shadows enabled for this side
	DynamicShadowsEnabled bool
}

// MapDispInfo
type MapDispInfo struct {
	// Face
	Face face.Face
	// EntityNum
	EntityNum int32
	// Power is power of this displacement (normally 2-4)
	Power int32
	// MinTess
	MinTess int32
	// SmoothingAngle
	SmoothingAngle float32
	// UAxis
	UAxis mgl32.Vec3
	// VAxis
	VAxis mgl32.Vec3
	// StartPosition
	StartPosition mgl32.Vec3
	// AlphaValues
	AlphaValues []float32
	// MaxDispDist
	MaxDispDist float32
	// DispDists
	DispDists []float32
	// VectorDisps
	VectorDisps []mgl32.Vec3
	// VectorOffsets
	VectorOffsets []mgl32.Vec3
	// Contents
	Contents int32
	// BrushSideID
	BrushSideID int32
	// Flags
	Flags int32

	// #ifdef VSVMFIO
	// Elevation float32
	// OffsetNormals []Vector
}

// CompressedLightCube
type CompressedLightCube struct {
	// Color
	Color [6]ColorRGBExponent32
}

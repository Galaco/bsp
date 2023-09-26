package common

import (
	"github.com/galaco/bsp/primitives/face"
	"github.com/go-gl/mathgl/mgl32"
)

// ColorRGBExponent32 is a 4 byte color representation (RGBExp)
type ColorRGBExponent32 struct {
	// R red
	R uint8 `json:"R"`
	// G green
	G uint8 `json:"G"`
	// B blue
	B uint8 `json:"B"`
	// Exponent exponent
	Exponent int8 `json:"Exponent"`
}

// Winding
type Winding struct {
	// Original
	Original int32 `json:"original"` // qboolean = int32
	// NumPoints
	NumPoints int32 `json:"numPoints"`
	// Points
	Points []mgl32.Vec3 `json:"points"`
}

// Side
type Side struct {
	// PlaneNum
	PlaneNum int32 `json:"planeNum"`
	// TexInfo
	TexInfo int32 `json:"texInfo"`
	// MapDisp
	MapDisp *MapDispInfo `json:"mapDisp"`
	// Winding
	Winding *Winding `json:"winding"`
	// Original
	Original *Side `json:"original"`
	// Contents
	Contents int32 `json:"contents"`
	// Surf
	Surf int32 `json:"surf"`
	// Visible
	Visible int32 `json:"visible"`
	// Tested
	Tested int32 `json:"tested"`
	// Bevel
	Bevel int32 `json:"bevel"`
	// Next
	Next *Side `json:"next"`
	// OrigIndex
	OrigIndex int32 `json:"origIndex"`
	// Id is side Id
	Id int32 `json:"id"`
	// SmoothingGroups
	SmoothingGroups uint32 `json:"smoothingGroups"`
	// AOverlayIds
	AOverlayIds []int32 `json:"aOverlayIds"`
	// AWaterOverlayIds
	AWaterOverlayIds []int32 `json:"aWaterOverlayIds"`
	// DynamicShadowsEnabled are dynamic shadows enabled for this side
	DynamicShadowsEnabled bool `json:"dynamicShadowsEnabled"`
}

// MapDispInfo
type MapDispInfo struct {
	// Face
	Face face.Face `json:"face"`
	// EntityNum
	EntityNum int32 `json:"entityNum"`
	// Power is power of this displacement (normally 2-4)
	Power int32 `json:"power"`
	// MinTess
	MinTess int32 `json:"minTess"`
	// SmoothingAngle
	SmoothingAngle float32 `json:"smoothingAngle"`
	// UAxis
	UAxis mgl32.Vec3 `json:"uAxis"`
	// VAxis
	VAxis mgl32.Vec3 `json:"vAxis"`
	// StartPosition
	StartPosition mgl32.Vec3 `json:"startPosition"`
	// AlphaValues
	AlphaValues []float32 `json:"alphaValues"`
	// MaxDispDist
	MaxDispDist float32 `json:"maxDispDist"`
	// DispDists
	DispDists []float32 `json:"dispDists"`
	// VectorDisps
	VectorDisps []mgl32.Vec3 `json:"vectorDisps"`
	// VectorOffsets
	VectorOffsets []mgl32.Vec3 `json:"vectorOffsets"`
	// Contents
	Contents int32 `json:"contents"`
	// BrushSideID
	BrushSideID int32 `json:"brushSideID"`
	// Flags
	Flags int32 `json:"flags"`

	// @TODO Figure out what this old C++ ifdef was for.
	// #ifdef VSVMFIO
	// Elevation float32
	// OffsetNormals []Vector
}

// CompressedLightCube
type CompressedLightCube struct {
	// Color
	Color [6]ColorRGBExponent32 `json:"color"`
}

package common

import (
	"github.com/galaco/bsp/primitives/face"
	"github.com/go-gl/mathgl/mgl32"
)

type ColorRGBExponent32 struct {
	R        uint8
	G        uint8
	B        uint8
	Exponent int8
}

const MAX_POINTS_ON_FIXED_WINDING = 12

//const MAX_DISPVERTS = 0
//const MAX_DISPTRIS = 0

type Winding struct {
	Original  int32 // qboolean = int32
	NumPoints int32
	Points    []mgl32.Vec3
}

type Side struct {
	PlaneNum              int32
	TexInfo               int32
	MapDisp               *MapDispInfo
	Winding               *Winding
	Original              *Side
	Contents              int32
	Surf                  int32
	Visible               int32
	Tested                int32
	Bevel                 int32
	Next                  *Side
	OrigIndex             int32
	Id                    int32
	SmoothingGroups       uint32
	AOverlayIds           []int32
	AWaterOverlayIds      []int32
	DynamicShadowsEnabled bool
}

type MapDispInfo struct {
	Face           face.Face
	EntityNum      int32
	Power          int32
	MinTess        int32
	SmoothingAngle float32
	uAxis          mgl32.Vec3
	vAxis          mgl32.Vec3
	StartPosition  mgl32.Vec3
	AlphaValues    []float32
	MaxDispDist    float32
	DispDists      []float32
	VectorDisps    []mgl32.Vec3
	VectorOffsets  []mgl32.Vec3
	Contents       int32
	BrushSideID    int32
	Flags          int32

	// #ifdef VSVMFIO
	// Elevation float32
	// OffsetNormals []Vector
}

package occlusion

import (
	"github.com/go-gl/mathgl/mgl32"
)

type OcclusionData struct {
	Flags int32
	FirstPoly int32
	PolyCount int32
	Mins mgl32.Vec3
	Maxs mgl32.Vec3
	Area int32
}

type OcclusionPolyData struct {
	FirstVertexIndex int32
	VertexCount int32
	PlaneNum int32
}

//Implementation info
/*
Lump size:
num occlusions * size()
+ num occlusionPolyDatas * size()
+ num OccluderVertexIndices * size()
+ 3* sizeof(Int32)
 */
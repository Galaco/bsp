package occlusion

import (
	"github.com/go-gl/mathgl/mgl32"
)

type OcclusionData struct {
	Flags     int32      `json:"flags"`
	FirstPoly int32      `json:"firstPoly"`
	PolyCount int32      `json:"polyCount"`
	Mins      mgl32.Vec3 `json:"mins"`
	Maxs      mgl32.Vec3 `json:"maxs"`
	Area      int32      `json:"area"`
}

type OcclusionPolyData struct {
	FirstVertexIndex int32 `json:"firstVertexIndex"`
	VertexCount      int32 `json:"vertexCount"`
	PlaneNum         int32 `json:"planeNum"`
}

//Implementation info
/*
Lump size:
num occlusions * size()
+ num occlusionPolyDatas * size()
+ num OccluderVertexIndices * size()
+ 3* sizeof(Int32)
*/

package lumps

import "github.com/go-gl/mathgl/mgl32"

// OcclusionData describes an occluder and its span in Occlusion.Polygons.
// HasArea records the presence of the area field in revision 2.
type OcclusionData struct {
	Flags     int32
	FirstPoly int32
	PolyCount int32
	Mins      mgl32.Vec3
	Maxs      mgl32.Vec3
	Area      int32
	HasArea   bool
}

// OcclusionPolyData describes a polygon and its span in Occlusion.VertexIndices.
// PlaneNum references the map's plane lump.
type OcclusionPolyData struct {
	FirstVertexIndex int32
	VertexCount      int32
	PlaneNum         int32
}

// Occlusion is the complete counted payload of lump 9. Wire counts are derived
// from these slices; vertex indices reference the map's vertex lump.
type Occlusion struct {
	Data          []OcclusionData
	Polygons      []OcclusionPolyData
	VertexIndices []int32
}

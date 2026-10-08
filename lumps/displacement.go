package lumps

import "github.com/go-gl/mathgl/mgl32"

const (
	// DisplacementFlagMagic marks the byte-24 word as flags instead of MinTess.
	DisplacementFlagMagic uint32 = 0x80000000
	// DisplacementFlagHasMultiblend is defined by the modern Reactive Drop dialect.
	DisplacementFlagHasMultiblend uint32 = 0x40000000
)

// Displacement exposes all 176 bytes of a Source displacement record.
// MinTess and Flags are mutually exclusive interpretations of the byte-24 word.
// Flags includes DisplacementFlagMagic and preserves every unknown flag bit.
// Inactive neighbor slots, padding and unused AllowedVerts bits are retained.
type Displacement struct {
	StartPosition       mgl32.Vec3
	DispVertStart       int32
	DispTriStart        int32
	Power               int32
	MinTess             int32
	SmoothingAngle      float32
	Contents            int32
	Flags               uint32
	MapFace             uint16
	MapFacePadding      [2]byte
	LightmapAlphaStart  int32
	LightmapSampleStart int32
	EdgeNeighbors       [4]DisplacementEdgeNeighbor
	CornerNeighbors     [4]DisplacementCornerNeighbor
	AllowedVerts        [10]uint32
}

// DisplacementSubNeighbor is six bytes; Index 0xffff makes all enum bytes inactive.
type DisplacementSubNeighbor struct {
	Index        uint16
	Orientation  uint8
	Span         uint8
	NeighborSpan uint8
	Padding      uint8
}

type DisplacementEdgeNeighbor struct {
	SubNeighbors [2]DisplacementSubNeighbor
}

// DisplacementCornerNeighbor retains inactive indices beyond Count.
type DisplacementCornerNeighbor struct {
	Neighbors [4]uint16
	Count     uint8
	Padding   uint8
}

// DisplacementSample identifies a local triangle and three quantized barycentric
// coordinates. The defined byte stream represents triangle indices 0 through 510.
type DisplacementSample struct {
	Index       uint16
	Barycentric [3]uint8
}

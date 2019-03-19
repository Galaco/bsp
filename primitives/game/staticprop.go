package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Staticprop lump contains all information
// relating to staticprop entries
type StaticPropLump struct {
	DictLump  StaticPropDictLump
	LeafLump  StaticPropLeafLump
	PropLumps []IStaticPropDataLump
}

// Flat array that consists of a unique list of all
// model filename+paths used by staticprops
type StaticPropDictLump struct {
	DictEntries int32
	Name        []string // Slice length must equal dictEntries. Validation to be added
}

// Represents a flat array of leaf indexes for all staticprops.
// A staticprop will have an offset and number on entries into the array that
// specify what leafs a given staticprop is contained in.
type StaticPropLeafLump struct {
	LeafEntries int32
	Leaf        []uint16 // Slice length must equal leafEntries. Validation to be added
}

// There are many different staticprop versions
// This interface should be up to date with all possible properties
// for any version.
// Missing properties across version should return 0,false,"" etc
type IStaticPropDataLump interface {
	GetOrigin() mgl32.Vec3
	GetAngles() mgl32.Vec3
	GetUniformScale() float32 // v11 onwards
	GetPropType() uint16
	GetFirstLeaf() uint16
	GetLeafCount() uint16
	GetSolid() uint8
	GetFlags() uint8
	GetSkin() int32
	GetFadeMinDist() float32
	GetFadeMaxDist() float32
	GetLightingOrigin() mgl32.Vec3
	GetForcedFadeScale() float32 // v5 onwards
	GetMinDXLevel() uint16       // v6+7 only
	GetMaxDXLevel() uint16       // v6+7 only
	GetMinCPULevel() uint8
	GetMaxCPULevel() uint8
	GetMinGPULevel() uint8
	GetMaxGPULevel() uint8
	GetDiffuseModulation() float32 // v7 onwards
	GetUnknown() float32           // v10 onwards
	GetDisableXBox360() bool       // v9 onwards
}

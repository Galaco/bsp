package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

// StaticPropLump contains all information
// relating to staticprop entries
type StaticPropLump struct {
	// DictLump
	DictLump StaticPropDictLump
	// LeafLump
	LeafLump StaticPropLeafLump
	// PropLumps
	PropLumps []IStaticPropDataLump
}

// StaticPropDictLump is a flat array that consists of a unique list of all
// model filename+paths used by staticprops
type StaticPropDictLump struct {
	// DictEntries
	DictEntries int32
	// Name
	Name []string // Slice length must equal dictEntries. Validation to be added
}

// StaticPropLeafLump represents a flat array of leaf indexes for all staticprops.
// A staticprop will have an offset and number on entries into the array that
// specify what leafs a given staticprop is contained in.
type StaticPropLeafLump struct {
	// LeafEntries
	LeafEntries int32
	// Leaf
	Leaf []uint16 // Slice length must equal leafEntries. Validation to be added
}

// IStaticPropDataLump There are many different staticprop versions
// This interface should be up to date with all possible properties
// for any version.
// Missing properties across version should return 0,false,"" etc
type IStaticPropDataLump interface {
	// GetOrigin Origin of object in world
	GetOrigin() mgl32.Vec3
	// GetAngles Rotation of object in world
	GetAngles() mgl32.Vec3
	// GetUniformScale Uniform scale of object in world
	// v11 onwards
	GetUniformScale() float32
	// GetPropType
	GetPropType() uint16
	// GetFirstLeaf Index into StaticPropLeafLump
	GetFirstLeaf() uint16
	// GetLeafCount Number of leafs this prop is in
	GetLeafCount() uint16
	// GetSolid
	GetSolid() uint8
	// GetFlags
	GetFlags() uint8
	// GetSkin Skin index of this prop
	GetSkin() int32
	// GetFadeMinDist
	GetFadeMinDist() float32
	// GetFadeMaxDist
	GetFadeMaxDist() float32
	// GetLightingOrigin World position to sample light from.
	GetLightingOrigin() mgl32.Vec3
	// GetForcedFadeScale
	// v5 onwards
	GetForcedFadeScale() float32
	// GetMinDXLevel Minimum directx level to render this prop
	// v6+7 only
	GetMinDXLevel() uint16
	// GetMaxDXLevel Maximum directx level to render this prop
	// v6+7 only
	GetMaxDXLevel() uint16
	// GetMinCPULevel Minimum CPU type to render this prop
	// v8 onwards
	GetMinCPULevel() uint8
	// GetMaxCPULevel Maximum CPU type to render this prop
	// v8 onwards
	GetMaxCPULevel() uint8
	// GetMinGPULevel
	// v8 onwards
	GetMinGPULevel() uint8
	// GetMaxGPULevel
	// v8 onwards
	GetMaxGPULevel() uint8
	// GetDiffuseModulation
	// v7 onwards
	GetDiffuseModulation() float32
	// GetUnknown
	// v10 onwards
	GetUnknown() float32
	// GetDisableXBox360 Should be disabled on xbox 360?
	// v9 onwards
	GetDisableXBox360() bool
}

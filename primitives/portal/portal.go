package portal

import (
	"github.com/galaco/bsp/primitives/common"
	"github.com/galaco/bsp/primitives/face"
	"github.com/galaco/bsp/primitives/node"
	"github.com/galaco/bsp/primitives/plane"
)

// MaxPortals is the maximum number of portals that can be generated
const MaxPortals = 65536

// Portal
type Portal struct {
	// Id portal id
	Id int32
	// Plane
	Plane plane.Plane
	// OnNode
	OnNode *node.Node
	// Nodes
	Nodes [2]*node.Node
	// Next
	Next [2]*Portal
	// Winding
	Winding *common.Winding
	// SideFound
	SideFound int32 //qboolean = int32
	// Side
	Side *common.Side
	// Face
	Face [2]*face.Face
}

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
	Id int32 `json:"id"`
	// Plane
	Plane plane.Plane `json:"plane"`
	// OnNode
	OnNode *node.Node `json:"onNode"`
	// Nodes
	Nodes [2]*node.Node `json:"nodes"`
	// Next
	Next [2]*Portal `json:"next"`
	// Winding
	Winding *common.Winding `json:"winding"`
	// SideFound
	SideFound int32 `json:"sideFound"` //qboolean = int32
	// Side
	Side *common.Side `json:"side"`
	// Face
	Face [2]*face.Face `json:"face"`
}

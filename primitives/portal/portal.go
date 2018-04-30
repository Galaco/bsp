package portal

import (
	"github.com/galaco/bsp/primitives/plane"
	"github.com/galaco/bsp/primitives/node"
	"github.com/galaco/bsp/primitives/face"
	"github.com/galaco/bsp/primitives/common"
)

const MAX_PORTALS = 65536

type Portal struct {
	Id int32
	Plane plane.Plane
	OnNode *node.Node
	Nodes [2]*node.Node
	Next [2]*Portal
	Winding *common.Winding
	SideFound int32 //qboolean = int32
	Side *common.Side
	Face [2]*face.Face
}

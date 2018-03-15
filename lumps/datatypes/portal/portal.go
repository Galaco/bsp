package portal

import (
	"github.com/galaco/bsp/lumps/datatypes/plane"
	"github.com/galaco/bsp/lumps/datatypes/node"
	"github.com/galaco/bsp/lumps/datatypes/face"
	"github.com/galaco/bsp/lumps/datatypes/common"
)

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

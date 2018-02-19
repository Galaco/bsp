package model

import "github.com/galaco/bsp/lumps/datatypes/common"

type Model struct {
	Mins common.Vector
	Maxs common.Vector
	Origin common.Vector
	HeadNode int32
	FirstFace int32
	NumFaces int32
}

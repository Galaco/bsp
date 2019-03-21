package leaf

import (
	"log"
	"testing"
)

func TestLeaf_Area(t *testing.T) {
	sut := getPrimitive()
	if sut.Area() != 61 {
		log.Println(sut.Area())
		t.Error("incorrect area value")
	}
}

func TestLeaf_Flags(t *testing.T) {
	sut := getPrimitive()
	if sut.Flags() != 46 {
		log.Println(sut.Flags())
		t.Error("incorrect flags value")
	}
}

func TestLeaf_SetArea(t *testing.T) {
	TestLeaf_Area(t)
}

func TestLeaf_SetFlags(t *testing.T) {
	TestLeaf_Flags(t)
}

func getPrimitive() Leaf {
	l := Leaf{}
	l.Contents = 1
	l.Cluster = 2
	//l.Name = 61 // ("a")[0]
	l.SetArea(61)
	l.SetFlags(46)
	l.Mins[0] = 3
	l.Mins[1] = 7
	l.Mins[2] = 14
	l.Maxs[0] = 13
	l.Maxs[1] = 36
	l.Maxs[2] = 51
	l.FirstLeafFace = 32
	l.NumLeafFaces = 11
	l.FirstLeafBrush = 21
	l.NumLeafBrushes = 76
	l.LeafWaterDataID = 621

	return l
}

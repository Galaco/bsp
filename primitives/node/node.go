package node

type Node struct {
	PlaneNum  int32
	Children  [2]int32
	Mins      [3]int16
	Maxs      [3]int16
	FirstFace uint16
	NumFaces  uint16
	Area      int16
	Padding   int16
}

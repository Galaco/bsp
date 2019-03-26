package node

// Node is a bsp node that has 2 children, because it was splittable along a plane.
// It represents a non-convex volume with child volumes
type Node struct {
	// PlaneNum is index of splitting place
	PlaneNum int32
	// Children are indexes of children that resulted from the split. It can reference either other Nodes or Leafs
	Children [2]int32
	// Mins is bounding volume min
	Mins [3]int16
	// Maxs is bounding volume maxs
	Maxs [3]int16
	// FirstFace is index into faces lump data
	FirstFace uint16
	// NumFaces is number of faces in this node
	NumFaces uint16
	// Area
	Area int16
	// Padding pads this structure to 4byte alignment (no other documented purpose)
	Padding int16
}

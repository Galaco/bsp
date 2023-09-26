package primitive

// Primitive
type Primitive struct {
	// Type
	Type byte `json:"type"`
	// FirstIndex
	FirstIndex uint16 `json:"firstIndex"`
	// IndexCount
	IndexCount uint16 `json:"indexCount"`
	// FirstVert
	FirstVert uint16 `json:"firstVert"`
	// VertCount
	VertCount uint16 `json:"vertCount"`
}

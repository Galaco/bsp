package physcollide

// PhysCollideEntry
type PhysCollideEntry struct {
	// ModelHeader
	ModelHeader PhysModel `json:"modelHeader"`
	// Solids
	Solids []Solid `json:"solids"`
	// TextBuffer
	TextBuffer string `json:"textBuffer"`
}

// PhysModel
type PhysModel struct {
	// ModelIndex
	ModelIndex int32 `json:"modelIndex"`
	// DataSize
	DataSize int32 `json:"dataSize"`
	// KeydataSize
	KeydataSize int32 `json:"keydataSize"`
	// SolidCount
	SolidCount int32 `json:"solidCount"`
}

// Solid
type Solid struct {
	// Size
	Size int32 `json:"size"`
	// CollisionBinary
	CollisionBinary []int32 `json:"collisionBinary"`
}

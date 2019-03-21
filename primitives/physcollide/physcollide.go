package physcollide

// PhysCollideEntry
type PhysCollideEntry struct {
	// ModelHeader
	ModelHeader PhysModel
	// Solids
	Solids []Solid
	// TextBuffer
	TextBuffer string
}

// PhysModel
type PhysModel struct {
	// ModelIndex
	ModelIndex int32
	// DataSize
	DataSize int32
	// KeydataSize
	KeydataSize int32
	// SolidCount
	SolidCount int32
}

// Solid
type Solid struct {
	// Size
	Size int32
	// CollisionBinary
	CollisionBinary []int32
}

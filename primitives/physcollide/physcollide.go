package physcollide

type PhysCollideEntry struct {
	ModelHeader PhysModel
	Solids []Solid
	TextBuffer string
}

type PhysModel struct {
	ModelIndex int32
	DataSize int32
	KeydataSize int32
	SolidCount int32
}

type Solid struct {
	Size int32
	CollisionBinary []int32
}
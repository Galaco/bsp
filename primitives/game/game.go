package game

// StaticPropLumpId is the identifier of the staticprop lump stored
const StaticPropLumpId = 1936749168

// Header
type Header struct {
	// LumpCount is the number of data lumps contained in the game lump
	LumpCount int32
	// GameLumps contains location and metadata about contained lumps
	// Slice length must equal lumpCount. Validation to be added
	GameLumps []LumpDef
}

// SetLumpCount set number of lumps
func (header *Header) SetLumpCount(num int32) *Header {
	header.LumpCount = num
	header.GameLumps = make([]LumpDef, header.LumpCount)
	return header
}

// LumpDef contains meta and location info about a lump contained within the game
// lump
type LumpDef struct {
	// Id is lump id
	Id int32
	// Flags is lump flags
	Flags uint16
	// Versionis lump version
	Version uint16
	// FileOffset is absolute offset into whole bsp
	FileOffset int32
	// FileLength is length of lump
	FileLength int32
}

// GenericGameLump represents a game lump with unknown/unmappable data
type GenericGameLump struct {
	// Data is byte representation of lump data
	Data []byte
}

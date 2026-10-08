package lumps

// StaticPropID is the Source game-directory identifier for static props (sprp).
const StaticPropID uint32 = 0x73707270

// GameLump is an owned snapshot of a game-directory entry, including unknown IDs.
// Offset is the encoded absolute input position, not a planned export position.
// StoredData is opaque stored bytes; compressed entries can have a larger LogicalLength.
// A final zero ID/flags/version/length entry is a preserved terminal boundary.
type GameLump struct {
	ID            uint32
	Flags         uint16
	Version       uint16
	Offset        int32
	LogicalLength int32
	StoredData    []byte
}

package lumps

// MapFlags is the map-wide flag word in lump 59. Unknown bits are preserved.
// Presence is separate from this value: a zero word can be explicitly present.
type MapFlags struct {
	LevelFlags uint32
}

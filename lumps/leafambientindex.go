package lumps

// AmbientIndex is a four-byte leaf ambient index from lump 51 or 52.
// For a zero sample count, FirstAmbientSample refers to another leaf's index,
// rather than to the sample array.
type AmbientIndex struct {
	AmbientSampleCount uint16
	FirstAmbientSample uint16
}

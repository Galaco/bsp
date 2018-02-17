package visibility


type Vis struct {
	NumClusters int32
	ByteOffset [][2]int32 // slice is incorrect. Slice length must equal numClusters (validation to be added)
}

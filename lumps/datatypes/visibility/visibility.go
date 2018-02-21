package visibility


type Vis struct {
	NumClusters int32
	ByteOffset [][2]int32 // Slice length = NumClusters [0]=offset to PVS bit array for cluster
	BitVectors []byte // Compressed bit vectors, contains run-length compression PVS data
}

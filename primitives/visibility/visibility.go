package visibility

const MAX_MAP_VISIBILITY= 0x1000000
const DVIS_PVS = 0
const DVIS_PAS = 1
const MAX_CLUSTER_SIZE_PER_VIS = 8

type Vis struct {
	NumClusters int32
	ByteOffset [][2]int32 // Slice length = NumClusters [0]=offset to PVS bit array for cluster
	BitVectors []byte // Compressed bit vectors, contains run-length compression PVS data
}

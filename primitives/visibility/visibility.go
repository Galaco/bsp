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

// Decompress vis data for a given cluster
// see https://developer.valvesoftware.com/wiki/Source_BSP_File_Format#Visibility for more
func (vis *Vis) GetPVSForCluster(clusterId int16) ([]bool) {
	visibleClusterIds := make([]bool, vis.NumClusters)
	v := vis.ByteOffset[clusterId][0] // pvs offset for cluster

	for c := int16(0); c < int16(vis.NumClusters); v++ {
		if vis.BitVectors[v] == 0 {
			v++
			c += int16(8 * uint8(vis.BitVectors[v]))
		} else {
			for bit := uint8(1); bit != 0; bit, c = bit * 2, c + 1 {
				if (vis.BitVectors[v] & bit) == 1 {
					visibleClusterIds[c] = true
				}
			}
		}
	}

	return visibleClusterIds
}
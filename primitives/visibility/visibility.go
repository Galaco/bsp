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

func (vis *Vis) GetVisibleClusters(clusterId int16) (visibleClusters []int16) {
	pvs := vis.GetPVSForCluster(clusterId)

	for id,visible := range pvs {
		if visible == true {
			visibleClusters = append(visibleClusters, int16(id))
		}
	}

	return visibleClusters
}

// Decompress vis data for a given cluster
// see https://developer.valvesoftware.com/wiki/Source_BSP_File_Format#Visibility for more
func (vis *Vis) GetPVSForCluster(clusterId int16) ([]bool) {
	visibleClusterIds := make([]bool, vis.NumClusters)
	v := vis.ByteOffset[clusterId][0] // pvs offset for cluster

	for c := int16(0); c < int16(vis.NumClusters); v++ {
		if uint8(vis.BitVectors[v]) == 0 {
			v++
			c += int16(8 * uint8(vis.BitVectors[v]))
		} else {
			for bit := uint8(1); bit != 0; bit, c = bit * 2, c + 1 {
				if c == int16(vis.NumClusters) {
					break
				}
				if (uint8(vis.BitVectors[v]) & bit) != 0 {
					visibleClusterIds[c] = true
				}
			}
		}
	}

	return visibleClusterIds
}
package visibility

const MAX_MAP_VISIBILITY = 0x1000000
const DVIS_PVS = 0
const DVIS_PAS = 1
const MAX_CLUSTER_SIZE_PER_VIS = 8

// Visibility data for clusters
// Includes both PVS (Potential Visible Set) and PAS (Potential Audible Set)
type Vis struct {
	NumClusters int32
	ByteOffset  [][2]int32 // Slice length = NumClusters [0]=offset to PVS bit array for cluster
	BitVectors  []byte     // Compressed bit vectors, contains run-length compression PVS data
}

// Get all visible clusters ids for a given cluster
func (vis *Vis) GetVisibleClusters(clusterId int16) (visibleClusters []int16) {
	pvs := vis.GetPVSForCluster(clusterId)

	for id, visible := range pvs {
		if visible == true {
			visibleClusters = append(visibleClusters, int16(id))
		}
	}

	return visibleClusters
}

// Decompress vis data for a given cluster
// see https://developer.valvesoftware.com/wiki/Source_BSP_File_Format#Visibility for more
func (vis *Vis) GetPVSForCluster(clusterId int16) []bool {
	visibleClusterIds := make([]bool, vis.NumClusters)
	v := vis.ByteOffset[clusterId][0] // pvs offset for cluster

	for currentClusterIdx := int32(0); currentClusterIdx < vis.NumClusters; v++ {
		if int(vis.BitVectors[v]) == 0 {
			v++
			currentClusterIdx += (int32(vis.BitVectors[v]) << 3)
			continue
		}
		for i := uint8(0); i < 8 && currentClusterIdx+int32(i) < vis.NumClusters; i++ {
			if (vis.BitVectors[v] & (1 << i)) != 0 {
				visibleClusterIds[currentClusterIdx+int32(i)] = true
			}
		}
		currentClusterIdx += 8
	}

	return visibleClusterIds
}

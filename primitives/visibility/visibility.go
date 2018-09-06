package visibility

const MAX_MAP_VISIBILITY = 0x1000000
const DVIS_PVS = 0
const DVIS_PAS = 1
const MAX_CLUSTER_SIZE_PER_VIS = 8

type Vis struct {
	NumClusters int32
	ByteOffset  [][2]int32 // Slice length = NumClusters [0]=offset to PVS bit array for cluster
	BitVectors  []byte     // Compressed bit vectors, contains run-length compression PVS data
}

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

	for offset := int32(0); offset < vis.NumClusters; v++ {
		if vis.BitVectors[v] == 0 {
			v++
			offset += int32(8 * vis.BitVectors[v])
			continue
		}
		for i := uint8(0); i < 8 && offset+int32(i) < vis.NumClusters; i++ {
			if (vis.BitVectors[v] & (1 << i)) != 0 {
				visibleClusterIds[offset+int32(i)] = true
			}
		}
		offset += 8
	}

	return visibleClusterIds
}

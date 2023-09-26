package visibility

// MaxMapVisibility is the maximum size of visibility data in bytes
const MaxMapVisibility = 0x1000000

// VisPVS index
const VisPVS = 0

// VisPAS index
const VisPAS = 1

// Vis contains visibility data for clusters
// Includes both PVS (Potential Visible Set) and PAS (Potential Audible Set)
type Vis struct {
	// NumClusters is number of computed cluster
	NumClusters int32 `json:"numClusters"`
	// ByteOffset contains offsets for cluster pvs and pas
	// Slice length = NumClusters [0]=offset to PVS bit array for cluster
	ByteOffset [][2]int32 `json:"byteOffset"`
	// BitVectors are compressed bit vectors, contains run-length compression PVS data
	BitVectors []byte `json:"bitVectors"`
}

// GetVisibleClusters returns all visible clusters ids for a given cluster
func (vis *Vis) GetVisibleClusters(clusterId int16) (visibleClusters []int16) {
	pvs := vis.GetPVSForCluster(clusterId)

	for id, visible := range pvs {
		if visible {
			visibleClusters = append(visibleClusters, int16(id))
		}
	}

	return visibleClusters
}

// GetPVSForCluster decompresses vis data for a given cluster
// see https://developer.valvesoftware.com/wiki/Source_BSP_File_Format#Visibility for more
func (vis *Vis) GetPVSForCluster(clusterId int16) []bool {
	visibleClusterIds := make([]bool, vis.NumClusters)
	v := vis.ByteOffset[clusterId][VisPVS] // pvs offset for cluster

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

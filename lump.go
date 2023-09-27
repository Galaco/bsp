package bsp

import (
	"fmt"

	"github.com/galaco/bsp/lumps"
	"github.com/galaco/bsp/versions"
)

// getReferenceLumpByIndex Return an instance of a Lump for a given offset.
func getReferenceLumpByIndex(index int, version int32) (lumps.Lump, error) {
	if index < 0 || index > 63 {
		return nil, fmt.Errorf("invalid lump id: %d provided", index)
	}

	l, err := versions.GetLumpForVersion(int(version), index)
	if err != nil {
		return nil, err
	}

	l.SetVersion(version)

	return l, nil
}

// DefaultLumpOrdering is Source Engines default export order (at least it's an order specified in the original SDK).
// Source compile tools write lumps out of order, so while the ordering doesn't actually matter, it may be
// useful/more performant to maintain the same order, particularly post-export.
var DefaultLumpOrdering = [64]LumpId{
	LumpPlanes,
	LumpLeafs,
	LumpLeafAmbientLighting,
	LumpLeafAmbientIndex,
	LumpLeafAmbientIndexHDR,
	LumpLeafAmbientLightingHDR,
	LumpVertexes,
	LumpNodes,
	LumpTexInfo,
	LumpTexData,
	LumpDispInfo,
	LumpDispVerts,
	LumpDispTris,
	LumpDispLightmapSamplePositions,
	LumpFaceMacroTextureInfo,
	LumpPrimitives,
	LumpPrimVerts,
	LumpPrimIndices,
	LumpFaces,
	LumpFacesHDR,
	LumpFaceIds,
	LumpOriginalFaces,
	LumpBrushes,
	LumpBrushSides,
	LumpLeafFaces,
	LumpLeafBrushes,
	LumpSurfEdges,
	LumpEdges,
	LumpModels,
	LumpAreas,
	LumpAreaPortals,
	LumpLighting,
	LumpLightingHDR,
	LumpVisibility,
	LumpEntities,
	LumpWorldLights,
	LumpWorldLightsHDR,
	LumpLeafWaterData,
	LumpOcclusion,
	LumpMapFlags,
	LumpPortals,
	LumpClusters,
	LumpPortalVerts,
	LumpClusterPortals,
	LumpClipPortalVerts,
	LumpCubemaps,
	LumpTexDataStringData,
	LumpTexDataStringTable,
	LumpOverlays,
	LumpWaterOverlays,
	LumpOverlayFades,
	LumpPhysCollide,
	LumpPhysDisp,
	LumpVertNormals,
	LumpVertNormalIndices,
	LumpLeafMinDistToWater,
	LumpGame,
	LumpPakfile,
}

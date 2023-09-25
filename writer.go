package bsp

import (
	"bytes"
	"encoding/binary"
	"sort"

	"github.com/galaco/bsp/lumps"
)

// Writer is a Bsp export writer.
type Writer struct {
}

// Write bsp to []byte.
func (w *Writer) Write(data *Bsp) ([]byte, error) {
	// First we need to update the header to reflect any lump changes
	// At the same time we can dump our lumps as bytes to write later
	lumpBytes := make([][]byte, 64)
	currentOffset := 1036 // Header always 1036bytes

	for _, index := range resolveLumpExportOrder(data.Header()) {
		// We have to handle lump 35 (GameData differently)
		// Because valve mis-designed the file format and relatively positioned data contains absolute file offsets.
		if index == LumpGame {
			gamelump := data.lumps[index].(*lumps.Game)
			data.lumps[index] = gamelump.UpdateInternalOffsets(int32(currentOffset) - data.header.Lumps[index].Offset)
		}
		exportBytes, err := data.Lump(index).ToBytes()
		if err != nil {
			return nil, err
		}
		lumpBytes[index] = exportBytes

		lumpSize := len(lumpBytes[index])

		data.header.Lumps[index].Length = int32(lumpSize)
		data.header.Lumps[index].Offset = int32(currentOffset)

		currentOffset += lumpSize

		// Finally 4byte align each lump.
		lumpBytes[index] = append(lumpBytes[index], make([]byte, currentOffset%4)...)
		currentOffset += currentOffset % 4
	}

	// Now we can export our bsp
	var buf bytes.Buffer

	//Write Header
	err := binary.Write(&buf, binary.LittleEndian, data.header)
	if err != nil {
		return nil, err
	}

	//Write lumps
	for _, lumpData := range lumpBytes {
		if err = binary.Write(&buf, binary.LittleEndian, lumpData); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func resolveLumpExportOrder(header *Header) [64]LumpId {
	vals := make([]struct {
		Offset int32
		Id     int
	}, 64)
	for idx, l := range header.Lumps {
		vals[idx].Offset = l.Offset
		vals[idx].Id = idx
	}

	sort.Slice(vals, func(i, j int) bool {
		return vals[i].Offset < vals[j].Offset
	})

	res := [64]LumpId{}

	for idx := range vals {
		res[idx] = LumpId(vals[idx].Id)
	}
	return res
}

// NewWriter Returns a new bsp writer instance.
func NewWriter() Writer {
	return Writer{}
}

// DefaultLumpOrdering is Source Engines default export order.
// Source compile tools write lumps out of order
// While the ordering doesn't actually matter, it may
// be useful/more performant to maintain the same order, particularly post-export
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

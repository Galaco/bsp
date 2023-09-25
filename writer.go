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
	marshalledLumps := make([][]byte, 64)
	var err error

	exportOrder := resolveLumpExportOrder(data.Header())
	currentOffset := 1032 // Header always 1032bytes
	for _, index := range exportOrder {
		// We have to handle lump 35 (GameData differently)
		// Because valve designed the file format oddly and relatively positioned data contains absolute file offsets.
		if index == LumpGame {
			gamelump := data.lumps[index].(*lumps.Game)
			data.lumps[index] = gamelump.UpdateInternalOffsets(int32(currentOffset) - data.header.Lumps[index].Offset)
		}

		// Export lump to bytes.
		if marshalledLumps[index], err = data.Lump(index).ToBytes(); err != nil {
			return nil, err
		}

		if len(marshalledLumps[index]) == 0 {
			// 0 bytes is a valid lump, but the offset is apparently 0 for this case (likely
			// to avoid the cost of seeking).
			data.header.Lumps[index].Offset = 0
			continue
		}

		// Update header with new lump offset and length.
		data.header.Lumps[index].Offset = int32(currentOffset)
		data.header.Lumps[index].Length = int32(len(marshalledLumps[index]))

		currentOffset += int(data.header.Lumps[index].Length)

		// Finally 4byte align each lump.
		marshalledLumps[index] = append(marshalledLumps[index], make([]byte, currentOffset%4)...)
		currentOffset += currentOffset % 4
	}

	// Now we can export our bsp.
	var buf bytes.Buffer

	// Write Header.
	if err := binary.Write(&buf, binary.LittleEndian, data.header); err != nil {
		return nil, err
	}

	// Write lumps.
	for _, lumpData := range marshalledLumps {
		if err := binary.Write(&buf, binary.LittleEndian, lumpData); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// resolveLumpExportOrder returns a lump export order that matches the order of the original file
// based on lump offsets.
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
func NewWriter() *Writer {
	return &Writer{}
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

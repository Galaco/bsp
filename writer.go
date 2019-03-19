package bsp

import (
	"bytes"
	"encoding/binary"
	"github.com/galaco/bsp/lumps"
)

// Writer is a Bsp export writer.
type Writer struct {
	data Bsp
}

// GetBsp Gets bsp file to write.
func (w *Writer) GetBsp() Bsp {
	return w.data
}

// SetBsp Sets bsp file to write.
func (w *Writer) SetBsp(file Bsp) {
	w.data = file
}

// Write bsp to []byte.
func (w *Writer) Write() ([]byte, error) {
	// First we need to update the header to reflect any lump changes
	// At the same time we can dump our lumps as bytes to write later
	lumpBytes := make([][]byte, 64)
	currentOffset := 1036 // Header always 1036bytes

	for _, index := range getDefaultLumpOrdering() {
		// We have to handle lump 35 (GameData differently)
		// Because valve mis-designed the file format and relatively positioned data contains absolute file offsets.
		if index == LumpGame {
			gamelump := w.data.lumps[int(index)].GetContents().(*lumps.Game)
			w.data.lumps[int(index)].SetContents(
				gamelump.UpdateInternalOffsets(int32(currentOffset) - w.data.header.Lumps[int(index)].Offset))
		}
		exportBytes, err := w.WriteLump(index)
		if err != nil {
			return nil, err
		}
		lumpBytes[int(index)] = exportBytes

		lumpSize := len(lumpBytes[int(index)])

		w.data.header.Lumps[int(index)].Length = int32(lumpSize)
		w.data.header.Lumps[int(index)].Offset = int32(currentOffset)

		currentOffset += lumpSize

		// Finally 4byte align each lump.
		lumpBytes[int(index)] = append(lumpBytes[int(index)], make([]byte, currentOffset%4)...)
		currentOffset += currentOffset % 4
	}

	// Now we can export our bsp
	var buf bytes.Buffer

	//Write Header
	err := binary.Write(&buf, binary.LittleEndian, w.data.header)
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

// WriteLump Exports a single lump to []byte.
func (w *Writer) WriteLump(index LumpId) ([]byte, error) {
	lump := w.data.GetLump(index)
	return lump.Marshall()
}

// NewWriter Returns a new bsp writer instance.
func NewWriter() Writer {
	w := Writer{}
	return w
}

// getDefaultLumpOrdering gets Source Engines default export order.
// Source compile tools write lumps out of order
// While the ordering doesn't actually matter, it may
// be useful/more performant to maintain the same order, particularly post-export
func getDefaultLumpOrdering() [64]LumpId {
	return [64]LumpId{
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
}

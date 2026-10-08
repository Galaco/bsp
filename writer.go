package bsp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

// Writer binds a file explicitly; export never mutates its authoritative payloads.
type Writer struct{ data *Bsp }

// Write bsp to []byte.
func (w *Writer) Write() ([]byte, error) {
	if w.data == nil {
		return nil, fmt.Errorf("%w: nil file", ErrMalformed)
	}
	header := w.data.header
	// First we need to update the header to reflect any lump changes
	// At the same time we can dump our lumps as bytes to write later
	lumpBytes := make([][]byte, 64)
	currentOffset := 1036 // Header always 1036bytes

	for _, index := range getDefaultLumpOrdering() {
		// Plan locally; never decode or mutate authoritative payloads during export.
		if currentOffset > math.MaxInt32 {
			return nil, fmt.Errorf("%w: output offset exceeds int32", ErrMalformed)
		}
		var exportBytes []byte
		var err error
		if index == LumpGame {
			exportBytes, err = w.data.MarshalGameLump(int32(currentOffset))
		} else {
			exportBytes, err = w.WriteLump(index)
		}
		if err != nil {
			return nil, err
		}
		lumpBytes[int(index)] = exportBytes

		lumpSize := len(lumpBytes[int(index)])

		header.Lumps[int(index)].Length = int32(lumpSize)
		header.Lumps[int(index)].Offset = int32(currentOffset)

		currentOffset += lumpSize

		// Finally 4byte align each lump.
		lumpBytes[int(index)] = append(lumpBytes[int(index)], make([]byte, currentOffset%4)...)
		currentOffset += currentOffset % 4
	}

	// Now we can export our bsp
	var buf bytes.Buffer

	//Write Header
	headerBytes := encodeHeader(header, w.data.layout)
	buf.Write(headerBytes[:])

	//Write lumps
	for _, lumpData := range lumpBytes {
		if err := binary.Write(&buf, binary.LittleEndian, lumpData); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// WriteLump Exports a single lump to []byte.
func (w *Writer) WriteLump(index LumpId) ([]byte, error) {
	if w.data == nil {
		return nil, fmt.Errorf("%w: nil file", ErrMalformed)
	}
	return w.data.RawLump(index)
}

// NewWriter binds the supplied file without a misleading value-copy facade.
func NewWriter(file *Bsp) *Writer { return &Writer{data: file} }

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

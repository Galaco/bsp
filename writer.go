package bsp


import (
	"bytes"
	"encoding/binary"
	"github.com/galaco/bsp/lumps"
)

// Bsp export writer.
type Writer struct {
	data Bsp
}

// Get bsp file to write.
func (w *Writer) GetBsp() Bsp {
	return w.data
}

// Set bsp file to write.
func (w *Writer) SetBsp(file Bsp) {
	w.data = file
}

// Write bsp to []byte.
func (w *Writer) Write() []byte {
	// First we need to update the header to reflect any lump changes
	// At the same time we can dump our lumps as bytes to write later
	lumpBytes := make([][]byte, 64)
	currentOffset := 1036 // Header always 1036bytes

	for _,index := range getDefaultLumpOrdering() {
		// We have to handle lump 35 (GameData differently)
		// Because valve mis-designed the file format and relatively positioned data contains absolute file offsets.
		if index == LUMP_GAME_LUMP {
			gamelump := (w.data.lumps[index].GetContents()).(lumps.Game)
			w.data.lumps[index].SetContents(
				gamelump.UpdateInternalOffsets(int32(currentOffset) - w.data.header.Lumps[index].Offset))
		}
		lumpBytes[index] = w.WriteLump(index)

		lumpSize := len(lumpBytes[index])

		w.data.header.Lumps[index].Length = int32(lumpSize)
		w.data.header.Lumps[index].Offset = int32(currentOffset)

		currentOffset += lumpSize

		// Finally 4byte align each lump.
		lumpBytes[index] = append(lumpBytes[index], make([]byte, currentOffset % 4)...)
		currentOffset += currentOffset % 4
	}

	// Now we can export our bsp
	var buf bytes.Buffer

	//Write Header
	binary.Write(&buf, binary.LittleEndian, w.data.header)

	//Write lumps
	for _,lumpData := range lumpBytes {
		binary.Write(&buf, binary.LittleEndian, lumpData)

	}

	return buf.Bytes()
}

/**
	Export a single lump to []byte.
 */
func (w *Writer) WriteLump(index int) []byte {
	lump := w.data.GetLump(index)
	return lump.GetContents().ToBytes()
}

/**
	Return a new bsp writer instance.
 */
func NewWriter() Writer {
	w := Writer{}
	return w
}

func getDefaultLumpOrdering() [64]int {
	return [64]int {
		LUMP_PLANES,
		LUMP_LEAFS,
		LUMP_LEAF_AMBIENT_LIGHTING,
		LUMP_LEAF_AMBIENT_INDEX,
		LUMP_LEAF_AMBIENT_INDEX_HDR,
		LUMP_LEAF_AMBIENT_LIGHTING_HDR,
		LUMP_VERTEXES,
		LUMP_NODES,
		LUMP_TEXINFO,
		LUMP_TEXDATA,
		LUMP_DISPINFO,
		LUMP_DISP_VERTS,
		LUMP_DISP_TRIS,
		LUMP_DISP_LIGHTMAP_SAMPLE_POSITIONS,
		LUMP_FACE_MACRO_TEXTURE_INFO,
		LUMP_PRIMITIVES,
		LUMP_PRIMVERTS,
		LUMP_PRIMINDICES,
		LUMP_FACES,
		LUMP_FACES_HDR,
		LUMP_FACEIDS,
		LUMP_ORIGINALFACES,
		LUMP_BRUSHES,
		LUMP_BRUSHSIDES,
		LUMP_LEAFFACES,
		LUMP_LEAFBRUSHES,
		LUMP_SURFEDGES,
		LUMP_EDGES,
		LUMP_MODELS,
		LUMP_AREAS,
		LUMP_AREAPORTALS,
		LUMP_LIGHTING,
		LUMP_LIGHTING_HDR,
		LUMP_VISIBILITY,
		LUMP_ENTITIES,
		LUMP_WORLDLIGHTS,
		LUMP_WORLDLIGHTS_HDR,
		LUMP_LEAFWATERDATA,
		LUMP_OCCLUSION,
		LUMP_MAP_FLAGS,
		LUMP_PORTALS,
		LUMP_CLUSTERS,
		LUMP_PORTALVERTS,
		LUMP_CLUSTERPORTALS,
		LUMP_CLIPPORTALVERTS,
		LUMP_CUBEMAPS,
		LUMP_TEXDATA_STRING_DATA,
		LUMP_TEXDATA_STRING_TABLE,
		LUMP_OVERLAYS,
		LUMP_WATEROVERLAYS,
		LUMP_OVERLAY_FADES,
		LUMP_PHYSCOLLIDE,
		LUMP_PHYSDISP,
		LUMP_VERTNORMALS,
		LUMP_VERTNORMALINDICES,
		LUMP_LEAFMINDISTTOWATER,
		LUMP_GAME_LUMP,
		LUMP_PAKFILE,
		}
}
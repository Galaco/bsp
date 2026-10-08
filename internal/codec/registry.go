package codec

import (
	"fmt"
	"github.com/galaco/bsp/internal/format"
)

// NewDecoder selects a codec using the actual lump revision and engine dialect.
func NewDecoder(index int, context format.Context) (Decoder, error) {
	if err := ValidateRevision(index, context); err != nil {
		return nil, err
	}
	var ret Decoder
	switch index {
	case 0:
		ret = new(EntData)
	case 1:
		ret = new(Planes)
	case 2:
		ret = new(TextureData)
	case 3:
		ret = new(Vertices)
	case 4:
		ret = new(Visibility)
	case 5:
		ret = new(Nodes)
	case 6:
		ret = new(TextureInfo)
	case 7:
		ret = new(Faces)
	case 8:
		ret = new(Lighting)
	case 9:
		ret = new(Occlusion)
	case 10:
		ret = new(Leaf)
	case 11:
		ret = new(FaceIDs)
	case 12:
		ret = new(Edges)
	case 13:
		ret = new(SurfaceEdges)
	case 14:
		ret = new(Models)
	case 15:
		ret = new(WorldLight)
	case 16:
		ret = new(LeafFaces)
	case 17:
		ret = new(LeafBrushes)
	case 18:
		ret = new(Brushes)
	case 19:
		ret = new(BrushSide)
	case 20:
		ret = new(Areas)
	case 21:
		ret = new(AreaPortals)
	case 22:
		ret = new(PropCollisions)
	case 23:
		ret = new(PropHulls)
	case 24:
		ret = new(PropHullVertices)
	case 25:
		ret = new(PropHullTriangles)
	case 26:
		ret = new(Displacements)
	case 27:
		ret = new(Faces)
	case 28, 29:
		return nil, fmt.Errorf("opaque lump %d", index)
	case 30:
		ret = new(VertexNormals)
	case 31:
		ret = new(VertexNormalIndices)
	case 32:
		return nil, fmt.Errorf("opaque lump %d", index)
	case 33:
		ret = new(DisplacementVertices)
	case 34, 35:
		return nil, fmt.Errorf("opaque lump %d", index)
	case 36:
		ret = new(LeafWaterData)
	case 37:
		ret = new(Primitives)
	case 38:
		ret = new(PrimitiveVertices)
	case 39:
		ret = new(PrimitiveIndices)
	case 40:
		ret = new(Pakfile)
	case 41:
		ret = new(ClipPortalVertices)
	case 42:
		ret = new(Cubemaps)
	case 43:
		ret = new(TextureStringData)
	case 44:
		ret = new(TextureStringTable)
	case 45:
		ret = new(Overlays)
	case 46:
		ret = new(LeafMinDistancesToWater)
	case 47:
		ret = new(FaceMacroTextureInfo)
	case 48:
		ret = new(DisplacementTriangles)
	case 49:
		return nil, fmt.Errorf("byte-stream lump %d; use attached PropBlob accessor", index)
	case 50:
		ret = new(WaterOverlays)
	case 51:
		ret = new(LeafAmbientIndexHDR)
	case 52:
		ret = new(LeafAmbientIndex)
	case 53:
		ret = new(Lighting) //lighting hdr
	case 54:
		ret = new(WorldLightHDR)
	case 55:
		ret = new(LeafAmbientLightingHDR)
	case 56:
		ret = new(LeafAmbientLighting)
	case 57:
		return nil, fmt.Errorf("opaque lump %d", index)
	case 58:
		ret = new(Faces)
	case 59:
		ret = new(MapFlags)
	case 60:
		ret = new(OverlayFades)
	case 61:
		ret = new(OverlaySystemLevels)
	case 62:
		return nil, fmt.Errorf("opaque lump %d", index)
	case 63:
		ret = new(DisplacementMultiblends)
	default:
		return nil, fmt.Errorf("invalid lump %d", index)
	}
	ret.SetContext(context)
	return ret, nil
}
